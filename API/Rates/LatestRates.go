package rates

import (
	"encoding/json"
	zlog "github.com/rs/zerolog/log"
	dataModules "goTestProj/DataModules"
	customError "goTestProj/Error"
	"net/http"
	"time"
)

//RespondWithLatestRates finds latest rates in DB and responds with them as JSON
func (db Database) RespondWithLatestRates(w http.ResponseWriter, r *http.Request) {
	rates, err := db.getLatestRates()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		zlog.Error().Err(err).Msg("Failed to get rates from database")
		json.NewEncoder(w).Encode(customError.JSONErrorResponse{Message: "Failed to get rates from database"})
		// trūkst return, izprintēsies arī nākamais write uz w.
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(rates)
}

func (db Database) getLatestRates() (dataModules.ResponseData, error) {
	var data dataModules.ResponseData

	newestDataDateString, err := db.getLatestEntryDateString()
	if err != nil {
		return data, err
	}

	rows, err := db.Database.Query("SELECT currency, rate FROM RATES WHERE pubdate=$1", newestDataDateString)

	if err != nil {
		return data, err
	}
	defer rows.Close()

	data.PubDate = newestDataDateString
	for rows.Next() {
		var rate dataModules.Rates
		if err := rows.Scan(&rate.Currency, &rate.Rate); err != nil {
			continue
		}
		data.Rates = append(data.Rates, rate)
	}

	return data, nil
}

func (db Database) getLatestEntryDateString() (string, error) {
	var response dataModules.ResponseData
	err := db.Database.QueryRow("SELECT max(pubdate) FROM rates").Scan(&response.PubDate)
	if err != nil {
		return "", err
	}
	layout := "2006-01-02T03:04:05Z"
	str := response.PubDate
	t, err := time.Parse(layout, str)

	if err != nil {
		return "", err
	}

	return t.Format("2006-01-02"), nil
}