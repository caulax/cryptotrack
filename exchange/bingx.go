package exchange

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const (
	bingxApiUrl = "https://open-api.bingx.com/openApi/spot/v1/ticker/price?symbol=%s-USDT" // BTC-USDT
)

type BingxTrade struct {
	Timestamp int64  `json:"timestamp"`
	TradeID   string `json:"tradeId"`
	Price     string `json:"price"`
	Amount    string `json:"amount"`
	Type      int    `json:"type"`
	Volume    string `json:"volume"`
}

type BingxSymbolData struct {
	Symbol string       `json:"symbol"`
	Trades []BingxTrade `json:"trades"`
}

type BingxApiResponse struct {
	Code      int               `json:"code"`
	Timestamp int64             `json:"timestamp"`
	Data      []BingxSymbolData `json:"data"`
}

func GetCoinPriceBingx(coinName string) float64 {
	// URL for the bingx API endpoint
	url := fmt.Sprintf(bingxApiUrl, coinName)

	// Create a new HTTP request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading HTTP response body:", err)
		os.Exit(1)
	}

	// Parse the JSON response
	var bingxResponse BingxApiResponse
	err = json.Unmarshal(body, &bingxResponse)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		os.Exit(1)
	}

	// Check if the response contains data
	if len(bingxResponse.Data) == 0 {
		fmt.Println(coinName, "No ticker data available.")
	}

	var returnValue float64
	// Get the latest price of coin
	for _, symbolData := range bingxResponse.Data {
		for _, trade := range symbolData.Trades {
			returnValue, _ = strconv.ParseFloat(trade.Price, 64)
		}
	}

	return returnValue
}
