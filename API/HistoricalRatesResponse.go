package api

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	customError "goTestProj/Error"
	"goTestProj/Repository"
	"net/http"
)

//DataSource holds pointer to repo for work with DB
type DataSource struct {
	Repo *repository.DataSource
}

//RespondWithHistoricalData checks if DB has data for requested currency and responds with JSON if found
func (dataSource DataSource) RespondWithHistoricalData(w http.ResponseWriter, r *http.Request) {

	currency := r.URL.Query().Get("currency")
	rates, err := dataSource.Repo.GetHistoricalData(currency)

	if err == customError.MissingRates() {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(customError.JSONErrorResponse{Message: "Failed to fetch rates from source"})
		return
	} else if err == customError.QueryError() {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(rates) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rates)
	} else {
		w.WriteHeader(http.StatusNotFound)
		errorMessage := "Could not find requested currency " + currency
		log.Error().Err(err).Str("method", r.Method).Str("url", r.RequestURI).Msg(errorMessage)
		json.NewEncoder(w).Encode(customError.JSONErrorResponse{Message: "Could not find rates requested currency"})
	}
}