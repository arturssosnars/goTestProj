package DataModules

type ResponseData struct {
	Rates	[]Rates  `json:"rates"`
	PubDate	string `json:"pubDate"`
}
