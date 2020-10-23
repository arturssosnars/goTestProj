package api

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	config "goTestProj/Config"
	"goTestProj/Repository"
	tasks "goTestProj/Tasks"
	"net/http"
)

//SetListeners sets router for API
func (dataSource DataSource) SetListeners() {
	router := chi.NewRouter()
	router.HandleFunc("/all", dataSource.RespondWithLatestRates)
	router.HandleFunc("/single", dataSource.RespondWithHistoricalData)

	taskDatabase := tasks.Database{Database: dataSource.Repo.Database}
	go taskDatabase.SetTaskForAddingRates()

	port := ":" + config.API().Port
	log.Log().Msg("Listen to port " + port + "...")
	err := http.ListenAndServe(port, router)
	log.Fatal().Err(err).Msg("Failed to launch API")
}

//InitializeDB initializes DB
func InitializeDB() DataSource {
	var db *sql.DB

	err := config.Init()
	if err != nil {
		log.Fatal().Err(err)
	}

	postgres := config.Postgres()
	url := fmt.Sprintf("%v:%v/%v", postgres.URL, postgres.Port, postgres.Database)

	pgURL, err := pq.ParseURL(url)
	log.Log().
		Str("url", postgres.URL).
		Str("port", postgres.Port).
		Str("db", postgres.Database).
		Str("driver", postgres.Driver).
		Msg("Parsing URL for connection to DB")

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse database URL")
	}

	db, err = sql.Open(postgres.Driver, pgURL)
	log.Log()

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open database")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to verify connection with database")
	}

	// Redzu, ka pirms servera uzstartēšanas, tu veic valūtas kursa vērtību sinhronizāciju, nevis realizē to kā
	// atsevišķu CLI komandu. Nekas, ar šo tiksim galā :)
	// Problēma ar šo pieeju ir tāda, ka, gadījumā, ja latvijas bankas rss feeds ir down, serveri nav iespējams uzstartēt.

	// Iesaku iepazīties ar https://github.com/spf13/cobra, un iznest atsevišķā komandā likmju importu uz DB
	saveDatabase := repository.DataSource{Database: db}
	err = saveDatabase.AddRatesToDB()
	if err != nil {
		log.Fatal().Err(err)
	}

	return DataSource{&repository.DataSource{Database: db}}
}