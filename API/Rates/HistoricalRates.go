package Rates

import (
	"encoding/json"
	zlog "github.com/rs/zerolog/log"
	"goTestProj/API/Initialization"
	dataModules "goTestProj/DataModules"
	customError "goTestProj/Error"
	"net/http"
	"time"
)

func RespondWithHistoricalData(w http.ResponseWriter, r *http.Request) {
	currency := r.URL.Query().Get("currency")

	if len(currency) != 3 {
		var errorMessage string
		if len(currency) == 0 {
			errorMessage = "Missing currency parameter"
		} else {
			errorMessage = "currency param length should be 3"
		}
		zlog.Error().Str("method", r.Method).Str("url", r.RequestURI).Msg(errorMessage)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(customError.ErrorResponse{"Could not find rates requested currency"})
		return
	}

	rates, err := getCurrencyHistoricalRates(currency)

	if err != nil {
		zlog.Error().Err(err).Msg("Failed to get rates from database")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(customError.ErrorResponse{"Failed to get rates from database"})
	}

	if len(rates) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rates)
	} else {
		w.WriteHeader(http.StatusNotFound)
		errorMessage := "Could not find requested currency " + currency
		zlog.Error().Err(err).Str("method", r.Method).Str("url", r.RequestURI).Msg(errorMessage)
		json.NewEncoder(w).Encode(customError.ErrorResponse{"Could not find rates requested currency"})
	}
}

func getCurrencyHistoricalRates(currency string) ([]dataModules.HistoricalRate, error) {
	var data []dataModules.HistoricalRate

	rows, err := Initialization.Db.Query("SELECT rate, pubdate FROM rates WHERE currency=$1", currency)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rate dataModules.HistoricalRate
		if err := rows.Scan(&rate.Rate, &rate.PubDate); err != nil {
			continue
		}
		layout := "2006-01-02T03:04:05Z"
		str := rate.PubDate
		t, err := time.Parse(layout, str)

		if err != nil {
			continue
		}

		rate.PubDate = t.Format("2006-01-02")
		data = append(data, rate)
	}

	return data, nil
}