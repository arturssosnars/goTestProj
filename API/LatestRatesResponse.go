package api

import (
	"encoding/json"
	zlog "github.com/rs/zerolog/log"
	customError "goTestProj/Error"
	"net/http"
)

//RespondWithLatestRates finds latest rates in DB and responds with them as JSON
func (dataSource DataSource) RespondWithLatestRates(w http.ResponseWriter, r *http.Request) {
	rates, err := dataSource.Repo.GetLatestRates()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		zlog.Error().Err(err).Msg("Failed to get rates from database")
		json.NewEncoder(w).Encode(customError.JSONErrorResponse{Message: "Failed to get rates from database"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(rates)
}