package API

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	dataModules "goTestProj/DataModules"
	"log"
	"net/http"
	"time"
)

func respondWithHistoricalData(w http.ResponseWriter, r *http.Request) {
	// 1. kļūdu variabļus ieteicams saukt par `err`, jo `error collides with builtin interface`
	// 2. Pastāv iespēja, ka handleris neradīs nekādu kļūdu, padarot šo mainīgo lieku šajā scopā.
	var error dataModules.Error

	// ja mums interesē tikai pirmā query parametra vērtība, var izmantot `r.URL.Query().Get("currency")`
	currency, ok := r.URL.Query()["currency"]

	if !ok || len(currency[0]) != 3 || len(currency) > 1 {
		error.Message = "Invalid parameters"
		RespondWithError(w, http.StatusBadRequest, error)
		return
	}

	// Šai funkcijai lielākais trūkums ir tāds, ka tā neatgriež kļūdu.
	// Jānis Kope integrēsies pret tavu API, paprasīs valūtas kursus, nesaņems atbildi, jo notikusi kļūda,
	// un paprasīs tev "Vecīt, kāpēc API man neko neatgriež? Kas notika?".
	// Lai atrisinātu šādas problēmas, vienmēr vajadzētu atgriezt kļūdu klientam, paskaidrojot kas noticis, un vajadzētu
	// logot kļūdu ar, iespējams, papildus detaļām, lai tu pats varētu atrast un saprast, kas nogāja greizi.
	rates := getCurrencyHistoricalRates(currency[0])

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(rates)
}

func getCurrencyHistoricalRates(currency string) []dataModules.HistoricalRate {
	var data []dataModules.HistoricalRate

	// padot query argumentus ar string concatenation prasās pēc SQL injekcijas
	// DB.Query ļauj norādīt vaicājumu ar placeholder vērtībām, kuras pēc tam pareizajā secībā var norādīt
	// kā varargs:
	// Db.Query("SELECT rate, pubdate FROM rates WHERE currency=$1", currency)
	rows, err := Db.Query("SELECT rate, pubdate FROM rates WHERE currency='" + currency + "';")

	// šajā gadījumā būtu labāk atgriezt kļūdu
	LogErrorIfNeeded(err)
	defer rows.Close()

	for rows.Next() {
		var rate dataModules.HistoricalRate
		// šajā gadījumā būtu labāk atgriezt kļūdu vai izlaist iterāciju
		if err := rows.Scan(&rate.Rate, &rate.PubDate); err != nil {
			log.Fatal(err)
		}
		layout := "2006-01-02T03:04:05Z"
		str := rate.PubDate
		t, err := time.Parse(layout, str)

		// šajā gadījumā būtu labāk atgriezt kļūdu vai izlaist iterāciju
		LogErrorIfNeeded(err)
		// šim uzdevumam šāds laika formāts der, jo granulārākā laika vienība ir diena, bet
		// citur es izmantotu vai nu ISO8601 datetime formātu, vai unix timestamp
		rate.PubDate = t.Format("2006-01-02")
		data = append(data, rate)
	}

	return data
}

// newest var arī uzrakstīt kā latest :)
func respondWithNewestRates(w http.ResponseWriter, r *http.Request) {
	rates := getNewestRates()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(rates)
}

func getNewestRates() dataModules.ResponseData {
	var data dataModules.ResponseData

	newestDataDateString := getNewestDate()
	// tas pats kas pirmīt par string concatenation
	rows, err := Db.Query("SELECT currency, rate FROM RATES WHERE pubdate='" + newestDataDateString + "';")

	// šajā gadījumā būtu labāk atgriezt kļūdu
	LogErrorIfNeeded(err)
	defer rows.Close()

	data.PubDate = newestDataDateString
	for rows.Next() {
		var rate dataModules.Rates
		// šajā gadījumā būtu labāk atgriezt kļūdu vai izliast iterāciju
		if err := rows.Scan(&rate.Currency, &rate.Rate); err != nil {
			log.Fatal(err)
		}
		data.Rates = append(data.Rates, rate)
	}

	return data
}

func getNewestDate() string {
	var response dataModules.ResponseData
	err := Db.QueryRow("SELECT max(pubdate) FROM rates").Scan(&response.PubDate)
	// šajā gadījumā būtu labāk atgriezt kļūdu
	LogErrorIfNeeded(err)
	layout := "2006-01-02T03:04:05Z"
	str := response.PubDate
	t, err := time.Parse(layout, str)

	// šajā gadījumā būtu labāk atgriezt kļūdu
	LogErrorIfNeeded(err)

	return t.Format("2006-01-02")
}

func GetBankRates() dataModules.Rss {
	// hardkodēts resursa parametrs
	resp, err := http.Get("https://www.bank.lv/vk/ecb_rss.xml")

	// šajā gadījumā būtu labāk atgriezt kļūdu
	LogErrorIfNeeded(err)

	var bankData = dataModules.Rss{}

	err = xml.NewDecoder(resp.Body).Decode(&bankData)

	// šajā gadījumā būtu labāk atgriezt kļūdu
	LogErrorIfNeeded(err)

	return bankData
}

func CreateQuery(rates dataModules.ResponseData) string {
	query := ""
	if !rowExists("SELECT id FROM rates WHERE pubDate=$1", rates.PubDate) {
		for rate := 0; rate < len(rates.Rates); rate++ {
			//workaround for low connection count (for free tier cloud solution)
			processedRate := rates.Rates[rate]
			statement := "INSERT INTO rates (currency, rate, pubDate) VALUES('" + processedRate.Currency + "'," + fmt.Sprintf("%f", processedRate.Rate) + ",'" + rates.PubDate + "');"
			query = query + statement
		}
	}

	return query
}

func rowExists(query string, args ...interface{}) bool {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := Db.QueryRow(query, args...).Scan(&exists)
	// šajā gadījumā būtu labāk atgriezt kļūdu
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	return exists
}