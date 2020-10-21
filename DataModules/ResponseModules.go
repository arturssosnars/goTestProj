package datamodules

//Rates holds info for single rate with currency and rate values
type Rates struct {
	Currency 	string 	`json:"currency"`
	Rate 		float64 `json:"rate"`
}

//HistoricalRate holds info for rate at single date
type HistoricalRate struct {
	Rate float64 `json:"rate"`
	PubDate string `json:"pubDate"`
}

//ResponseData is used to form JSON for API response when latest rates are requested
type ResponseData struct {
	Rates	[]Rates  `json:"rates"`
	PubDate	string `json:"pubDate"`
}