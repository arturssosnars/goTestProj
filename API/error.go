package API

import (
	"encoding/json"
	"goTestProj/DataModules"
	"log"
	"net/http"
)

func LogErrorIfNeeded(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func RespondWithError(w http.ResponseWriter, status int, error DataModules.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}