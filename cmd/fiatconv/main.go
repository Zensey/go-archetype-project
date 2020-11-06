package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var version string

type ratesDict map[string]float64

type exchangeRatesAPIResp struct {
	Rates ratesDict
	Base  string
	Date  string
}

func getExchangeRates(data *exchangeRatesAPIResp) error {
	res, err := http.Get("https://api.exchangeratesapi.io/latest")
	if err != nil {
		return err
	}
	return json.NewDecoder(res.Body).Decode(&data)
}

func getRatio(data exchangeRatesAPIResp, symbolSrc, symbolDst string) (float64, error) {
	rateSrc, ok := data.Rates[symbolSrc]
	if symbolSrc == data.Base {
		rateSrc = 1
	} else if !ok {
		return 0, fmt.Errorf("Invalid source currency: %s", symbolSrc)
	}

	rateDst, ok := data.Rates[symbolDst]
	if symbolDst == data.Base {
		rateDst = 1
	} else if !ok {
		return 0, fmt.Errorf("Invalid destination currency: %s", symbolDst)
	}

	return rateDst / rateSrc, nil
}

func main() {
	if len(os.Args[1:]) < 3 {
		fatal("Usage: fiatconv <amount_src:float> <src_symbol:string> <dst_symbol:string>\n\n" +
			"Arguments:\n" +
			"  amount_src  Amount to convert\n" +
			"  src_symbol  Currency you are converting from\n" +
			"  dst_symbol  Currency you are converting to\n\n" +
			"Example:\n" +
			"  fiatconv 100 EUR GBP")
	}

	amountSrc, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid amountSrc")
	}
	symbolSrc := os.Args[2]
	symbolDst := os.Args[3]

	data := exchangeRatesAPIResp{}
	err = getExchangeRates(&data)
	if err != nil {
		fatal("Request error: %v", err)
	}

	ratio, err := getRatio(data, symbolSrc, symbolDst)
	if err != nil {
		fatal("%v", err)
	}

	amountDst := ratio * amountSrc
	fmt.Printf("%.2f %s = %.2f %s\n", amountSrc, symbolSrc, amountDst, symbolDst)
	return
}

func fatal(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}
