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
	bybitApiUrl = "https://api.bybit.com/v5/market/tickers?category=spot&symbol=%sUSDT" // BTCUSDT
)

type BybitTickerResponse struct {
	RetCode    int        `json:"retCode"`
	RetMsg     string     `json:"retMsg"`
	Result     Result     `json:"result"`
	RetExtInfo RetExtInfo `json:"retExtInfo"`
	Time       int64      `json:"time"`
}

type Result struct {
	Category string       `json:"category"`
	List     []TickerData `json:"list"`
}

type TickerData struct {
	Symbol       string `json:"symbol"`
	Bid1Price    string `json:"bid1Price"`
	Bid1Size     string `json:"bid1Size"`
	Ask1Price    string `json:"ask1Price"`
	Ask1Size     string `json:"ask1Size"`
	LastPrice    string `json:"lastPrice"`
	PrevPrice24h string `json:"prevPrice24h"`
	Price24hPcnt string `json:"price24hPcnt"`
	HighPrice24h string `json:"highPrice24h"`
	LowPrice24h  string `json:"lowPrice24h"`
	Turnover24h  string `json:"turnover24h"`
	Volume24h    string `json:"volume24h"`
}

type RetExtInfo struct {
	// Add fields here if needed when they appear in the response
}

func GetCoinPriceBybit(coinName string) float64 {
	// URL for the bybit API endpoint
	url := fmt.Sprintf(bybitApiUrl, coinName)

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
	var bybitResponse BybitTickerResponse
	err = json.Unmarshal(body, &bybitResponse)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		os.Exit(1)
	}
	var bid1Price string
	// Check if the response contains data
	if len(bybitResponse.Result.List) > 0 {
		// Access the bid1Price field from the first element of the list
		bid1Price = bybitResponse.Result.List[0].Bid1Price
		fmt.Printf("The bid1Price for %s is: %s\n", bybitResponse.Result.List[0].Symbol, bid1Price)
	} else {
		fmt.Println("No ticker data available.")
	}

	var returnValue float64
	// Get the latest price of coin
	returnValue, _ = strconv.ParseFloat(bid1Price, 64)

	return returnValue
}
