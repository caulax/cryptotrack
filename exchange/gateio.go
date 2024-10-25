package exchange

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	gateioApiUrl = "https://api.gateio.ws/api/v4/spot/tickers?currency_pair=%s_USDT" // BTC_USDT
)

// Struct to map the JSON response from Gate.io
type GateioTickerData struct {
	CurrencyPair string `json:"currency_pair"`
	Last         string `json:"last"`
	LowestAsk    string `json:"lowest_ask"`
	HighestBid   string `json:"highest_bid"`
	Change       string `json:"change_percentage"`
	High24h      string `json:"high_24h"`
	Low24h       string `json:"low_24h"`
	Vol24h       string `json:"base_volume"`
	QuoteVol24h  string `json:"quote_volume"`
}

func GetCoinPriceGateio(coinName string) float64 {
	url := fmt.Sprintf(gateioApiUrl, coinName)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		os.Exit(1)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to fetch data:", err)
		os.Exit(1)
	}

	var ticker []GateioTickerData
	if err := json.NewDecoder(resp.Body).Decode(&ticker); err != nil {
		fmt.Println("Failed to decode data:", err)
		os.Exit(1)
	}

	var coinPrice float64
	if len(ticker) > 0 {
		coinPrice, _ = strconv.ParseFloat(ticker[0].Last, 64)
	}

	return coinPrice
}
