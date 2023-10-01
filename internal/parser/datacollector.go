package parser

import (
	"context"
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"

	"github.com/ArturKKK/CurrencyTransfer/internal/db"
	"github.com/ArturKKK/CurrencyTransfer/pkg/logging"
	"golang.org/x/net/html/charset"
)

func Parse(url string, db *db.Database, logger *logging.Logger) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Errorf("failed to create request: %v", err)
		return err
	}
	req.Header.Set("User-Agent", "Mygithub.com/ArturKKK/CurrencyTransferService/1.0")
	response, err := client.Do(req)
	if err != nil {
		logger.Errorf("failed to send request: %v", err)
		return err
	}

	defer response.Body.Close()

	var valCurs ValCurs
	decoder := xml.NewDecoder(response.Body)
	decoder.CharsetReader = charset.NewReaderLabel // windows-1251 to utf-8
	err = decoder.Decode(&valCurs)
	if err != nil {
		logger.Errorf("failed to decode response: %v", err)
		return err
	}

	for _, valute := range valCurs.Valutes {
		VunitRateStr := strings.Replace(valute.VunitRate, ",", ".", -1)
		vunitRate, err := strconv.ParseFloat(VunitRateStr, 64)
		if err != nil {
			logger.Errorf("failed to convert string to float64: %v", err)
			return err
		}
		db.Save(context.TODO(), valute.CharCode, vunitRate)
	}

	return nil
}
