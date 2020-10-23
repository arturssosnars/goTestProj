package repository

import "database/sql"

//Rates is used to form JSON response for latest rates
type Rates struct {
	Currency 	string 	`json:"currency"`
	Rate 		float64 `json:"rate"`
}

//ResponseData is used to form JSON for API response when
//latest rates are requested
type ResponseData struct {
	Rates	[]Rates  `json:"rates"`
	PubDate	string `json:"pubDate"`
}

//HistoricalRate holds info for rate at single date
type HistoricalRate struct {
	Rate float64 `json:"rate"`
	PubDate string `json:"pubDate"`
}

//DataSource is used to store pointer to data source
type DataSource struct {
	Database *sql.DB
}