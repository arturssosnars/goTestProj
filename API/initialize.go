package API

import (
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/lib/pq"
	"log"
	"net/http"
	"fmt"
)

// globāli mainīgie nav laba prakse, bet to noteikti zināji :)
// Iezīme A1: tā vietā, lai rakstītu http.Handler funkcijas, kas paļaujas uz globālu mainīgo,
// tu varētu izveidot structu, kas satur fieldu ar *sql.DB, un realizē visas nepieciešamās http.HandlerFunc
// kā receiver metodes.
var Db *sql.DB

func Initialize() {
	// Iztēlosimies programma logo šādas rindas:
	// `init db`
	// <log error>
	// <program dies>
	// Bieži vien problēma ir sliktu konfigurācijas parametru padošana servisam, tādēļ būtu labi arī izlogot,
	// uz kādu hostu, kādu portu, un kādu datubāzi programma cenšas savienoties, lai varētu pārliecināties, ka
	// problēma nav ar programmu vai tās konfigurāciju.
	// Logošanai iesaku lietot zerolog: https://github.com/rs/zerolog
	fmt.Println("init db")

	// Hardkodēt resursu parametrus nav laba prakse, bet to noteikti zināji :)
	// Šim nolūkam ieteicu viper bibliotēku: https://github.com/spf13/viper
	pgUrl, err := pq.ParseURL("postgres://pxjcqhji:MPY070OmuT4xUosjTGTWaHGH3jHCqJOY@hattie.db.elephantsql.com:5432/pxjcqhji")

	// skatīt komentārus pie funkcijas definīcijas
	LogErrorIfNeeded(err)

	Db, err = sql.Open("postgres", pgUrl)
	LogErrorIfNeeded(err)

	err = Db.Ping()
	LogErrorIfNeeded(err)

	// Redzu, ka pirms servera uzstartēšanas, tu veic valūtas kursa vērtību sinhronizāciju, nevis realizē to kā
	// atsevišķu CLI komandu. Nekas, ar šo tiksim galā :)
	// Problēma ar šo pieeju ir tāda, ka, gadījumā, ja latvijas bankas rss feeds ir down, serveri nav iespējams uzstartēt.
	AddRatesToDB()

	router := chi.NewRouter()
	router.HandleFunc("/all", respondWithNewestRates)
	router.HandleFunc("/single", respondWithHistoricalData)

	log.Println("set timers for tasks")

	// skatīt komentāru pie funkcijas
	go setTaskForAddingRates()

	// hardkodēts ports
	log.Println("Listen to port :8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}