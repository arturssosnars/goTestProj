package repository

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/go-ping/ping"
	zlog "github.com/rs/zerolog/log"
	config "goTestProj/Config"
	customError "goTestProj/Error"
	"net/http"
)

// Šis strukts ir mulsinošs. strukts ar nosaukumu Database, un ar pointeri uz DB konekciju liek man domāt,
// ka šeit ir definētas visas data-access funkcijas, bet nē, te ir implementēts valūtas kursu imports.
// Iesaku padomāt labāk pie nosaukumiem.

// Package savedata zina, uz kurieni vajag veikt kādu pieprasījumu, lai saņemtu kaut kādu XML,
// bet package dataModules zina, kāds formāts būs šim saņemtajam XML. Manuprāt visām zināšanām par RSS barotni
// būtu jābūt vienā pakotnē.
func getBankRates() (Rss, error) {
	var bankData Rss
	bankURL := config.Bank().URL
	resp, err := http.Get(bankURL)

	if err != nil {
		pinger, err2 := ping.NewPinger(bankURL)
		if err2 != nil {
			zlog.Error().Err(err2)
		} else {
			pinger.Count = 5
			pinger.OnFinish = func(stats *ping.Statistics) {
				if stats.PacketLoss > 50 {
					zlog.Fatal().Msg(fmt.Sprintf("Bad or no connection to bank API. Ping result: %v percent lost", stats.PacketLoss))
				}
			}
			pinger.Run()
		}

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
func (db DataSource) AddRatesToDB() error {
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

func (db DataSource) createQuery(rates ResponseData) (string, error) {
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

func (db DataSource) rowExists(query string, args ...interface{}) (bool, error) {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.Database.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}