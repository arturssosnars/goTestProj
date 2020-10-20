package SaveData

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	zlog "github.com/rs/zerolog/log"
	"goTestProj/API/Initialization"
	dataModules "goTestProj/DataModules"
	customError "goTestProj/Error"
	"net/http"
)

func GetBankRates() (dataModules.Rss, error) {
	var bankData dataModules.Rss
	// hardkodēts resursa parametrs
	resp, err := http.Get("https://www.bank.lv/vk/ecb_rss.xml")

	if err != nil {
		err = customError.BankApiError()
		return bankData, err
	}

	err = xml.NewDecoder(resp.Body).Decode(&bankData)

	if err != nil {
		err := customError.ParsingError()
		return bankData, err
	}

	return bankData, nil
}

func AddRatesToDB() {
	bankData, err := GetBankRates()

	if err != nil {
		if err == customError.ParsingError() {
			zlog.Error().Err(err).Msg("Failed to parse XML from bank API response")
		} else {
			zlog.Error().Err(err)
		}
	}

	rates, err := bankData.LatestRates()

	if err != nil {
		if err == customError.MissingRates() {
			zlog.Error().Err(err)
		} else if err == customError.ParsingError() {
			zlog.Error().Err(err).Msg("Failed to parse string array into ResponseData struct")
		}
	}

	query, err := CreateQuery(rates)

	if err != nil {
		zlog.Error().Err(err).Msg("Failed to create query")
	}

	_, err = Initialization.Db.Exec(query)

	if err != nil {
		zlog.Error().Err(err).Msg("Failed to insert new rates into database")
	}
}

func CreateQuery(rates dataModules.ResponseData) (string, error) {
	var query string
	rowExists, err := rowExists("SELECT id FROM rates WHERE pubDate=$1", rates.PubDate)
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

func rowExists(query string, args ...interface{}) (bool, error) {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := Initialization.Db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}