package repository

import (
	"time"
)

//GetLatestRates returns latest rates from DB
func (db DataSource) GetLatestRates() (ResponseData, error) {
	var data ResponseData

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
		var rate Rates
		if err := rows.Scan(&rate.Currency, &rate.Rate); err != nil {
			continue
		}
		data.Rates = append(data.Rates, rate)
	}

	return data, nil
}

func (db DataSource) getLatestEntryDateString() (string, error) {
	var response ResponseData
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