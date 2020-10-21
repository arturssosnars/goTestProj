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

	// config.Init() atgriež kļūdu, ja kaut kas neizdodas.
	config.Init()
	postgres := config.Postgres()
	// stringus var formatēt ar fmt.Sprintf
	url := postgres.URL + ":" + postgres.Port + "/" + postgres.Database

	// Tā vietā, lai katrā kļūdas paziņojumā duplicētu loga detaļas par savienojumu, kādu DB centās izveidot, labāk varētu
	// vienreiz logot, ka programma mēģina izveidot savienojumu ar tādiem un šādiem parametriem, un otro reizi, ka savienošanās
	// ar DB neizdevās dēļ šādiem iemesliem.
	pgURL, err := pq.ParseURL(url)

	// šis bloks neaptur programmau rīkoties tālāk ar sliktu pgURL vērtību
	if err != nil {
		zlog.Error().Err(err).
			Str("url", postgres.URL).
			Str("port", postgres.Port).
			Str("db", postgres.Database).
			Msg("Failed to parse database URL")
	}

	db, err = sql.Open(postgres.Driver, pgURL)

	// šis bloks neaptur programmu rīkoties tālāk ar sliktu db konekcijas vērtību
	if err != nil {
		zlog.Error().Err(err).
			Str("driverName", postgres.Driver).
			Str("url", postgres.URL).
			Str("port", postgres.Port).
			Str("db", postgres.Database).
			Msg("Failed to open database")
	}

	err = db.Ping()
	// šis bloks neaptur programmu rīkoties tālāk ar nesasniedzamu DB
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

	// Iesaku iepazīties ar https://github.com/spf13/cobra, un iznest atsevišķā komandā likmju importu uz DB
	saveDatabase := savedata.Database{Database: db}
	// šī metode atgriež kļūdu, kas nav handlota
	saveDatabase.AddRatesToDB()

	respondDatabase := rates.Database{Database: db}
	router := chi.NewRouter()
	router.HandleFunc("/all", respondDatabase.RespondWithLatestRates)
	router.HandleFunc("/single", respondDatabase.RespondWithHistoricalData)

	// nevajadzētu izmantot 2 dažādus loggerus.
	log.Println("set timers for tasks")

	taskDatabase := tasks.Database{Database: db}
	go taskDatabase.SetTaskForAddingRates()

	port := ":" + config.API().Port
	log.Println("Listen to port " + port + "...")
	log.Fatal(http.ListenAndServe(port, router))
}