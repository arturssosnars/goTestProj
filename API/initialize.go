package API

import (
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/lib/pq"
	"log"
	"net/http"
	"fmt"
)

var Db *sql.DB

func Initialize() {
	fmt.Println("init db")
	pgUrl, err := pq.ParseURL("postgres://pxjcqhji:MPY070OmuT4xUosjTGTWaHGH3jHCqJOY@hattie.db.elephantsql.com:5432/pxjcqhji")
	LogErrorIfNeeded(err)

	Db, err = sql.Open("postgres", pgUrl)
	LogErrorIfNeeded(err)

	err = Db.Ping()
	LogErrorIfNeeded(err)

	AddRatesToDB()

	router := chi.NewRouter()
	router.HandleFunc("/all", respondWithNewestRates)
	router.HandleFunc("/single", respondWithHistoricalData)

	log.Println("set timers for tasks")
	go setTaskForAddingRates()

	log.Println("Listen to port :8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}