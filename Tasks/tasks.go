package tasks

import (
	"database/sql"
	"goTestProj/Repository"
	"time"
)

//Database is used to pass pointer to DB as receiver
type Database struct {
	Database *sql.DB
}

//calculates time 1 minute past midnight UTC or 15 minutes from now
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

//SetTaskForAddingRates sets rates addition task depending with different timing
//depending on success or fail of previous try
func (db Database) SetTaskForAddingRates() {
	ticker := time.NewTicker(getTickerTime(false))

	for {
		select {
		case <-ticker.C:
			repo := repository.DataSource{Database: db.Database}
			err := repo.AddRatesToDB()
			if err != nil {
				ticker.Reset(getTickerTime(true))
			} else {
				ticker.Reset(getTickerTime(false))
			}
		}
	}
}