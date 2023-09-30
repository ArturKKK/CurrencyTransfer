package main

import (
	"context"
	"encoding/xml"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ArturKKK/CurrencyTransfer/internal/db"

	"golang.org/x/net/html/charset"
)

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	CharCode  string `xml:"CharCode"`
	VunitRate string `xml:"VunitRate"`
}

func Parse(url string, db *db.Database) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mygithub.com/ArturKKK/CurrencyTransferService/1.0")
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	var valCurs ValCurs
	decoder := xml.NewDecoder(response.Body)
	decoder.CharsetReader = charset.NewReaderLabel // windows-1251 to utf-8
	err = decoder.Decode(&valCurs)
	if err != nil {
		log.Fatal(err)
	}

	for _, valute := range valCurs.Valutes {
		VunitRateStr := strings.Replace(valute.VunitRate, ",", ".", -1)
		vunitRate, err := strconv.ParseFloat(VunitRateStr, 64)
		if err != nil {
			log.Fatal(err)
		}
		db.Save(context.TODO(), valute.CharCode, vunitRate)
	}

}
