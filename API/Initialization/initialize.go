package initialization

import (
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/lib/pq"
	zlog "github.com/rs/zerolog/log"
	"goTestProj/API/Rates"
	"goTestProj/API/SaveData"
	"goTestProj/Config"
	"goTestProj/Tasks"
	"log"
	"net/http"
)

//Initialize connects to DB, gets latest currency rates, inserts them into DB and opens :8000 for listening
func Initialize() {
	var db *sql.DB

	config.Init()
	postgres := config.Postgres()
	url := postgres.URL + ":" + postgres.Port + "/" + postgres.Database

	pgURL, err := pq.ParseURL(url)

	// skatīt komentārus pie funkcijas definīcijas
	if err != nil {
		zlog.Error().Err(err).
			Str("url", postgres.URL).
			Str("port", postgres.Port).
			Str("db", postgres.Database).
			Msg("Failed to parse database URL")
	}

	db, err = sql.Open(postgres.Driver, pgURL)

	if err != nil {
		zlog.Error().Err(err).
			Str("driverName", postgres.Driver).
			Str("url", postgres.URL).
			Str("port", postgres.Port).
			Str("db", postgres.Database).
			Msg("Failed to open database")
	}

	err = db.Ping()
	if err != nil {
		zlog.Error().Err(err).
			Str("driverName", postgres.Driver).
			Str("url", postgres.URL).
			Str("port", postgres.Port).
			Str("db", postgres.Database).
			Msg("Failed to verify connection with database")
	}

	// Redzu, ka pirms servera uzstartēšanas, tu veic valūtas kursa vērtību sinhronizāciju, nevis realizē to kā
	// atsevišķu CLI komandu. Nekas, ar šo tiksim galā :)
	// Problēma ar šo pieeju ir tāda, ka, gadījumā, ja latvijas bankas rss feeds ir down, serveri nav iespējams uzstartēt.
	saveDatabase := savedata.Database{Database: db}
	saveDatabase.AddRatesToDB()

	respondDatabase := rates.Database{Database: db}
	router := chi.NewRouter()
	router.HandleFunc("/all", respondDatabase.RespondWithLatestRates)
	router.HandleFunc("/single", respondDatabase.RespondWithHistoricalData)

	log.Println("set timers for tasks")

	taskDatabase := tasks.Database{Database: db}
	go taskDatabase.SetTaskForAddingRates()

	port := ":" + config.API().Port
	log.Println("Listen to port " + port + "...")
	log.Fatal(http.ListenAndServe(port, router))
}