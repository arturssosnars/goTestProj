package DataModules

type Rates struct {
	Currency 	string 	`json:"currency"`
	Rate 		float64 `json:"rate"`
}

type HistoricalRate struct {
	Rate float64 `json:"rate"`
	PubDate string `json:"pubDate"`
}

type ResponseData struct {
	Rates	[]Rates  `json:"rates"`
	PubDate	string `json:"pubDate"`
}