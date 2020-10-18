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
	var error dataModules.Error
	currency, ok := r.URL.Query()["currency"]

	if !ok || len(currency[0]) != 3 || len(currency) > 1 {
		error.Message = "Invalid parameters"
		RespondWithError(w, http.StatusBadRequest, error)
		return
	}

	rates := getCurrencyHistoricalRates(currency[0])

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(rates)
}

func getCurrencyHistoricalRates(currency string) []dataModules.HistoricalRate {
	var data []dataModules.HistoricalRate
	rows, err := Db.Query("SELECT rate, pubdate FROM rates WHERE currency='" + currency + "';")
	LogErrorIfNeeded(err)
	defer rows.Close()

	for rows.Next() {
		var rate dataModules.HistoricalRate
		if err := rows.Scan(&rate.Rate, &rate.PubDate); err != nil {
			log.Fatal(err)
		}
		layout := "2006-01-02T03:04:05Z"
		str := rate.PubDate
		t, err := time.Parse(layout, str)

		LogErrorIfNeeded(err)
		rate.PubDate = t.Format("2006-01-02")
		data = append(data, rate)
	}

	return data
}

func respondWithNewestRates(w http.ResponseWriter, r *http.Request) {
	rates := getNewestRates()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(rates)
}

func getNewestRates() dataModules.ResponseData {
	var data dataModules.ResponseData

	newestDataDateString := getNewestDate()
	rows, err := Db.Query("SELECT currency, rate FROM RATES WHERE pubdate='" + newestDataDateString + "';")

	LogErrorIfNeeded(err)
	defer rows.Close()

	data.PubDate = newestDataDateString
	for rows.Next() {
		var rate dataModules.Rates
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
	LogErrorIfNeeded(err)
	layout := "2006-01-02T03:04:05Z"
	str := response.PubDate
	t, err := time.Parse(layout, str)

	LogErrorIfNeeded(err)

	return t.Format("2006-01-02")
}

func GetBankRates() dataModules.Rss {
	resp, err := http.Get("https://www.bank.lv/vk/ecb_rss.xml")

	LogErrorIfNeeded(err)

	var bankData = dataModules.Rss{}

	err = xml.NewDecoder(resp.Body).Decode(&bankData)

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
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	return exists
}