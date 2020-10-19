package DataModules

import "encoding/xml"

// 1. Mums neinteresē visi lauki, kas tiek atgriezti šajā XML, tikai publikācijas datums un kursu strings.
// 2. Mums interesējošai apakšstruktūrai `Item` var definēt atsevišķu struktu.
//
// type Item struct {
//     RawRates string `xml:"description"`
//     PubDate string `xml:"pubDate"`
// }
//
// Ieguvums būtu tāds, ka mēs varētu definēt metodes uz šī strukta, kas atvieglotu ar šo struktu saistītas
// datu apstrādes darbības, kas vairākās vietās tiek veiktas.
//
// func (i Item) PubDateTime() (time.Time, error) {
//     t, err := time.Parse(i.PubDate)
//     if err != nil {
//         return nil, err
//     }
//     return t, nil
// }
type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Atom    string   `xml:"atom,attr"`
	Channel struct {
		Text        string `xml:",chardata"`
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Link        struct {
			Text string `xml:",chardata"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
			Href string `xml:"href,attr"`
		} `xml:"link"`
		LastBuildDate string `xml:"lastBuildDate"`
		Generator     string `xml:"generator"`
		Image         struct {
			Text   string `xml:",chardata"`
			URL    string `xml:"url"`
			Title  string `xml:"title"`
			Link   string `xml:"link"`
			Width  string `xml:"width"`
			Height string `xml:"height"`
		} `xml:"image"`
		Language string `xml:"language"`
		Ttl      string `xml:"ttl"`
		Item     []struct {
			Text  string `xml:",chardata"`
			Title string `xml:"title"`
			Link  string `xml:"link"`
			Guid  struct {
				Text        string `xml:",chardata"`
				IsPermaLink string `xml:"isPermaLink,attr"`
			} `xml:"guid"`
			Description string `xml:"description"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}