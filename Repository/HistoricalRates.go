package repository

import (
	zlog "github.com/rs/zerolog/log"
	customError "goTestProj/Error"
	"time"
)

// Šis un LatestRates fails nujau ir labāks - handleri neizmanto kaut kādu globālu DB variabli, bet izmanto to
// kā dependency.
// Bet, varētu vēl labāk - rates.Database ir saplūduši 2 concerni - HTTP servēšana un DB loģika. Šī programma varētu
// izskatīties tā, ka vienā package ir Repository strukts, kas realizē visus vaicājumus sevī saturošajai datubāzei ar
// metodēm, kuras var izsaukt citas packages, kurām tas interesē. HTTP handlošanas loģiku var realizēt citā package, un
// kā dependency izmantot šo Repository package.

//GetHistoricalData returns historical rates for selected currency from DB
func (db DataSource)GetHistoricalData(currency string) ([]HistoricalRate, error) {
	if len(currency) != 3 {
		var errorMessage string
		if len(currency) == 0 {
			errorMessage = "Missing currency parameter"
		} else {
			errorMessage = "currency param length should be 3"
		}
		zlog.Error().Err(customError.MissingRates()).Msg(errorMessage)
		return nil, customError.MissingRates()
	}

	rates, err := db.getCurrencyHistoricalRates(currency)

	if err != nil {
		zlog.Error().Err(err).Msg("Failed to get rates from database")
		return nil, err
	}

	return rates, nil
}

func (db DataSource) getCurrencyHistoricalRates(currency string) ([]HistoricalRate, error) {
	var data []HistoricalRate

	rows, err := db.Database.Query("SELECT rate, pubdate FROM rates WHERE currency=$1", currency)

	if err != nil {
		return nil, customError.QueryError()
	}
	defer rows.Close()

	for rows.Next() {
		var rate HistoricalRate
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
