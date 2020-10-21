package savedata

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	zlog "github.com/rs/zerolog/log"
	config "goTestProj/Config"
	dataModules "goTestProj/DataModules"
	customError "goTestProj/Error"
	"net/http"
)

//Database is used to pass pointer to DB as receiver
type Database struct {
	Database *sql.DB
}

func getBankRates() (dataModules.Rss, error) {
	var bankData dataModules.Rss
	bankURL := config.Bank().URL
	resp, err := http.Get(bankURL)

	if err != nil {
		err = customError.BankAPIError()
		return bankData, err
	}

	err = xml.NewDecoder(resp.Body).Decode(&bankData)

	if err != nil {
		err := customError.ParsingError()
		return bankData, err
	}

	return bankData, nil
}

//AddRatesToDB requests latest rates, checks if DB has those rates already, if not - adds them to DB
func (db Database) AddRatesToDB() error {
	bankData, err := getBankRates()

	if err != nil {
		if err == customError.ParsingError() {
			zlog.Error().Err(err).Msg("Failed to parse XML from bank API response")
			return err
		}
		zlog.Error().Err(err)
		return err
	}

	rates, err := bankData.LatestRates()

	if err != nil {
		if err == customError.MissingRates() {
			zlog.Error().Err(err)
			return err
		} else if err == customError.ParsingError() {
			zlog.Error().Err(err).Msg("Failed to parse string array into ResponseData struct")
			return err
		}
	}

	query, err := db.createQuery(rates)

	if err != nil {
		zlog.Error().Err(err).Msg("Failed to create query")
		return err
	}

	_, err = db.Database.Exec(query)

	if err != nil {
		zlog.Error().Err(err).Msg("Failed to insert new rates into database")
		return err
	}

	return nil
}

func (db Database) createQuery(rates dataModules.ResponseData) (string, error) {
	var query string
	rowExists, err := db.rowExists("SELECT id FROM rates WHERE pubDate=$1", rates.PubDate)
	if err != nil {
		return query, err
	}
	if !rowExists {
		for rate := 0; rate < len(rates.Rates); rate++ {
			//workaround for low connection count (for free tier cloud solution)
			processedRate := rates.Rates[rate]
			statement := "INSERT INTO rates (currency, rate, pubDate) VALUES('" + processedRate.Currency + "'," + fmt.Sprintf("%f", processedRate.Rate) + ",'" + rates.PubDate + "');"
			query = query + statement
		}
	}

	return query, nil
}

func (db Database) rowExists(query string, args ...interface{}) (bool, error) {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.Database.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}