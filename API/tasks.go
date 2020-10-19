package API

import (
	"time"
)

// šī funkcija imo nav nepieciešama, jo tās rezultāts vienmēr būs 24h 1m
func getTickerTime() time.Duration {
	currentTime := time.Now()
	//1 minute after bank pub date
	startTime := currentTime.Truncate(24 * time.Hour).Add(24 * time.Hour + time.Minute)

	return startTime.Sub(currentTime)
}

func setTaskForAddingRates() {
	// šeit varēji inicializēt tickeri ar argumentu 24*time.Hour + time.Minute
	ticker := time.NewTicker(getTickerTime()).C

	for {
		select {
		case <-ticker:
			// Šis tickers turpinās tikšķēt ad infinitum ar norādīto laika periodu,
			// tas nav metams ārā pēc viena tikšķa.
			ticker = time.NewTicker(getTickerTime()).C
			AddRatesToDB()
		}
	}
}