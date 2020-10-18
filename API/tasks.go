package API

import (
	"time"
)

func getTickerTime() time.Duration {
	currentTime := time.Now()
	//1 minute after bank pub date
	startTime := currentTime.Truncate(24 * time.Hour).Add(24 * time.Hour + time.Minute)

	return startTime.Sub(currentTime)
}

func setTaskForAddingRates() {
	ticker := time.NewTicker(getTickerTime()).C

	for {
		select {
		case <-ticker:
			ticker = time.NewTicker(getTickerTime()).C
			AddRatesToDB()
		}
	}
}