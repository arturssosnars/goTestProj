package Tasks

import (
	"database/sql"
	"goTestProj/API/SaveData"
	"time"
)

type Database struct {
	Database *sql.DB
}

// šī funkcija imo nav nepieciešama, jo tās rezultāts vienmēr būs 24h 1m
func getTickerTime(retry bool) time.Duration {
	currentTime := time.Now()
	var startTime time.Time
	if retry {
		//15 minutes later
		startTime = currentTime.Add(15 * time.Minute)
	} else {
		//1 minute after midnight in UTC
		startTime = currentTime.Truncate(time.Hour).Add(24 * time.Hour + time.Minute)
	}

	return startTime.Sub(currentTime)
}

func (db Database) SetTaskForAddingRates() {
	ticker := time.NewTicker(getTickerTime(false))

	for {
		select {
		case <-ticker.C:
			// Šis tickers turpinās tikšķēt ad infinitum ar norādīto laika periodu,
			// tas nav metams ārā pēc viena tikšķa.
			var database = SaveData.Database{db.Database}
			err := database.AddRatesToDB()
			if err != nil {
				ticker.Reset(getTickerTime(true))
			} else {
				ticker.Reset(getTickerTime(false))
			}
		}
	}
}