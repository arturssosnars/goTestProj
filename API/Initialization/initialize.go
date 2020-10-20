package Initialization

import (
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/lib/pq"
	zlog "github.com/rs/zerolog/log"
	rates "goTestProj/API/Rates"
	"goTestProj/API/SaveData"
	"goTestProj/Tasks"
	"log"
	"net/http"
	config "goTestProj/Config"
)

func Initialize() {
	// Iztēlosimies programma logo šādas rindas:
	// `init db`
	// <log error>
	// <program dies>
	// Bieži vien problēma ir sliktu konfigurācijas parametru padošana servisam, tādēļ būtu labi arī izlogot,
	// uz kādu hostu, kādu portu, un kādu datubāzi programma cenšas savienoties, lai varētu pārliecināties, ka
	// problēma nav ar programmu vai tās konfigurāciju.
	// Logošanai iesaku lietot zerolog: https://github.com/rs/zerolog

	// Hardkodēt resursu parametrus nav laba prakse, bet to noteikti zināji :)
	// Šim nolūkam ieteicu viper bibliotēku: https://github.com/spf13/viper
	var db *sql.DB

	config.Init()
	postgres := config.Postgres()
	url := postgres.Url + ":" + postgres.Port + "/" + postgres.Database

	pgUrl, err := pq.ParseURL(url)

	// skatīt komentārus pie funkcijas definīcijas
	if err != nil {
		zlog.Error().Err(err).
			Str("url", postgres.Url).
			Str("port", postgres.Port).
			Str("db", postgres.Database).
			Msg("Failed to parse database URL")
	}

	db, err = sql.Open(postgres.Driver, pgUrl)

	if err != nil {
		zlog.Error().Err(err).
			Str("driverName", postgres.Driver).
			Str("url", postgres.Url).
			Str("port", postgres.Port).
			Str("db", postgres.Database).
			Msg("Failed to open database")
	}

	err = db.Ping()
	if err != nil {
		zlog.Error().Err(err).
			Str("driverName", postgres.Driver).
			Str("url", postgres.Url).
			Str("port", postgres.Port).
			Str("db", postgres.Database).
			Msg("Failed to verify connection with database")
	}

	// Redzu, ka pirms servera uzstartēšanas, tu veic valūtas kursa vērtību sinhronizāciju, nevis realizē to kā
	// atsevišķu CLI komandu. Nekas, ar šo tiksim galā :)
	// Problēma ar šo pieeju ir tāda, ka, gadījumā, ja latvijas bankas rss feeds ir down, serveri nav iespējams uzstartēt.
	saveDatabase := SaveData.Database{db}
	saveDatabase.AddRatesToDB()

	respondDatabase := rates.Database{db}
	router := chi.NewRouter()
	router.HandleFunc("/all", respondDatabase.RespondWithLatestRates)
	router.HandleFunc("/single", respondDatabase.RespondWithHistoricalData)

	log.Println("set timers for tasks")

	taskDatabase := Tasks.Database{db}
	go taskDatabase.SetTaskForAddingRates()

	port := ":" + config.Api().Port
	log.Println("Listen to port " + port + "...")
	log.Fatal(http.ListenAndServe(port, router))
}