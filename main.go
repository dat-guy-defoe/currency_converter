package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Currencies struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Date       string     `xml:"Date,attr"`
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	XMLName  xml.Name `xml:"Valute"`
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Nominal  string   `xml:"Nominal"`
	Name     string   `xml:"Name"`
	Value    string   `xml:"Value"`
}

func getHttpResponse(dateString string) *http.Response {
	os.Setenv("HTTP_PROXY", "")
	resp, err := http.Get("https://www.cbr-xml-daily.ru/daily_utf8.xml?date_req=" + dateString)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		layout := "02/01/2006"
		date, _ := time.Parse(layout, dateString)
		date = date.AddDate(0, 0, -1)
		dateString = date.Format("02/01/2006")

		resp = getHttpResponse(dateString)
	}

	return resp
}

func main() {
	var requiredCurrencies []string
	var date string
	args := os.Args[1:]

	if len(args) > 0 {
		requiredCurrencies = strings.Split(args[0], "/")
	} else {
		fmt.Println("NO ARGUMENTS RECEIVED")
		return
	}

	if len(args) == 2 {
		date = args[1]
	} else {
		currentTime := time.Now()
		date = currentTime.Format("02/01/2006")
	}

	resp := getHttpResponse(date)

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var currencies Currencies

	xml.Unmarshal(body, &currencies)

	searchMap := map[string]string{}
	for _, v := range requiredCurrencies {
		searchMap[v] = "1"
	}

	t, _ := time.Parse("02/01/2006", date)
	formattedDate := t.Format("2006-01-02")

	for i := 0; i < len(currencies.Currencies); i++ {
		currencyCode := currencies.Currencies[i].CharCode
		if found := searchMap[currencyCode]; found == "1" {
			fmt.Printf("%v	%v	%v	%v\n", formattedDate, currencyCode, strings.ReplaceAll(currencies.Currencies[i].Value, ",", "."), currencies.Currencies[i].Nominal)
		}
	}
	fmt.Printf("%v	%v	%v	%v\n", formattedDate, "RUB", "1", "1")
}
