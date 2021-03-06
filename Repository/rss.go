package repository

// vajag sakārtot importu secību, vismaz goland tā saka
import (
	"encoding/xml"
	customError "goTestProj/Error"
	"strconv"
	"strings"
	"time"
)

//Rss is used to parse XML from bank API
type Rss struct {
	XMLName xml.Name `xml:"rss"`
	// ja Channel ir pilnīgi lieks, un nepieciešams tikai, lai tiktu klāt pie Item,
	//tad anotācijā var definēt []Item ar anotāciju `xml:"channel>item"`
	Channel struct {
		Item []Item `xml:"item"`
	} `xml:"channel"`
}

//Item is used to parse XML from bank API
//It holds rates for single day
type Item struct {
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

//LatestRates returns JSON with latest rates
func (rss Rss) LatestRates() (ResponseData, error) {
	var data ResponseData
	index := len(rss.Channel.Item) - 1
	if rss.Channel.Item[index].Description == "" {
		err := customError.MissingRates()
		return data, err
	}

	item := rss.Channel.Item[index]

	rawRatesArray := strings.Split(item.Description, " ")

	data = parseSliceToStruct(rawRatesArray[:len(rawRatesArray)-1], item)
	if len(data.Rates) == 0 {
		err := customError.ParsingError()
		return data, err
	}
	return data, nil
}

func (i Item) getPubDate() (time.Time, error) {
	var pubTime time.Time
	layout := "Mon, 02 Jan 2006 03:04:05 -0700"
	pubTime, err := time.Parse(layout, i.PubDate)
	if err != nil {
		return pubTime, err
	}
	return pubTime, nil
}

func (rss Rss) getLatestItem() Item {
	index := len(rss.Channel.Item) - 1
	return rss.Channel.Item[index]
}

func stringSliceToResponseData(slice []string) []Rates {
	var rates []Rates
	for i := 0; i < len(slice); i += 2 {
		if floatRate, err := strconv.ParseFloat(slice[i+1], 64); err == nil {
			rates = append(rates, Rates{
				Currency: slice[i],
				Rate:     floatRate,
			})
		}
		continue
	}
	return rates
}

func parseSliceToStruct(rawRatesArray []string, item Item) ResponseData {
	rates := ResponseData{}
	pubDate, err := item.getPubDate()
	if err != nil {
		rates.PubDate = item.PubDate
	} else {
		date := pubDate.Format("2006-01-02")
		rates.PubDate = date
	}

	rates.Rates = stringSliceToResponseData(rawRatesArray)

	return rates
}