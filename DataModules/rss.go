package DataModules

import (
	"encoding/xml"
	"strconv"
	"strings"
	"time"
	customError "goTestProj/Error"
)

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

func (rss Rss) LatestRates() (ResponseData, error) {
	var data ResponseData
	index := len(rss.Channel.Item) - 1
	if rss.Channel.Item[index].Description == "" {
		err := customError.MissingRates()
		return data, err
	}

	// būtu jauki, ja mēs varētu vienkārši paprasīt bankData.Channel.Item.Rates(),
	// un tas atgrieztu slice ar kursu vērtību structiem
	rawRatesArray := strings.Split(rss.Channel.Item[index].Description, " ")

	data = parseXmlToStruct(rawRatesArray[:len(rawRatesArray)-1], rss, index)
	if len(data.Rates) == 0 {
		err := customError.ParsingError()
		return data, err
	}
	return data, nil
}

// Tā vietā, lai veiktu darbības ar RSS struktu šeit, var definēt receiver metodes, lai mēs pašam struktam varētu pajautāt
// Rates un formatētu PubDate.
// Skatīt komentārus rss.go failā
func parseXmlToStruct(rawRatesArray []string, bankData Rss, index int) ResponseData {
	rates := ResponseData{}
	layout := "Mon, 02 Jan 2006 03:04:05 -0700"
	str := bankData.Channel.Item[index].PubDate
	t, err := time.Parse(layout, str)

	if err != nil {
		rates.PubDate = str
	} else {
		date := t.Format("2006-01-02")
		rates.PubDate = date
	}

	for i := 0; i < (len(rawRatesArray) - 2); i += 2 {
		if floatRate, err := strconv.ParseFloat(rawRatesArray[i+1], 64); err == nil {
			rates.Rates = append(rates.Rates, Rates{
				Currency: rawRatesArray[i],
				Rate:     floatRate,
			})
		}
		continue
	}

	return rates
}