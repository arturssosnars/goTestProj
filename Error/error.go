package error

import (
	"errors"
)

var errAPI = errors.New("Failed to get rates from bank API")
var errParse = errors.New("Failed to parse data")
var errMissingRates = errors.New("Rates are missing in RSS struct")

// Kāpēc ir nepieciešamas speciālas funkcijas, kas atgriež kļūdas?
// Labāk padari kļūdas publiskas, un lieto tās, kur nepieciešamas.

//MissingRates returns error for missing rates related issues
func MissingRates() error {
	return errMissingRates
}

//BankAPIError returns error for bank API related issues
func BankAPIError() error {
	return errAPI
}

//ParsingError returns error for parsing related issues
func ParsingError() error {
	return errParse
}

//JSONErrorResponse is used to form JSON to client with error message
type JSONErrorResponse struct {
	Message string `json:"message"`
}