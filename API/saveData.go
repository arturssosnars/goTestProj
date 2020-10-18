package API

import (
	dataModules "goTestProj/DataModules"
	"log"
	"strconv"
	"strings"
	"time"
)

func AddRatesToDB() {
	var error dataModules.Error

	bankData := GetBankRates()

	index := len(bankData.Channel.Item) - 1
	if bankData.Channel.Item[index].Description == "" {
		error.Message = "Rates are missing"
		log.Fatal(error)
	}

	rawRatesArray := strings.Split(bankData.Channel.Item[index].Description, " ")

	rates := parseXmlToStruct(rawRatesArray[:len(rawRatesArray)-1], bankData, index)

	query := CreateQuery(rates)

	_, err := Db.Exec(query)

	LogErrorIfNeeded(err)
}

func parseXmlToStruct(rawRatesArray []string, bankData dataModules.Rss, index int) dataModules.ResponseData {
	rates := dataModules.ResponseData{}
	layout := "Mon, 02 Jan 2006 03:04:05 -0700"
	str := bankData.Channel.Item[index].PubDate
	t, err := time.Parse(layout, str)

	LogErrorIfNeeded(err)

	date := t.Format("2006-01-02")
	rates.PubDate = date

	for i := 0; i < (len(rawRatesArray) - 2); i += 2 {
		if floatRate, err := strconv.ParseFloat(rawRatesArray[i+1], 64); err == nil {
			rates.Rates = append(rates.Rates, dataModules.Rates{
				Currency: rawRatesArray[i],
				Rate:     floatRate,
			})
		}
		LogErrorIfNeeded(err)
	}

	return rates
}