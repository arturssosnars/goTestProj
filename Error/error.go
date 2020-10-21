package Error

import (
	"errors"
)

var apiError = errors.New("Failed to get rates from bank API")
var parseError = errors.New("Failed to parse data")
var missingRates = errors.New("Rates are missing in RSS struct")

func MissingRates() error {
	return missingRates
}

func BankApiError() error {
	return apiError
}

func ParsingError() error {
	return parseError
}

type ErrorResponse struct {
	Message string `json:"message"`
}