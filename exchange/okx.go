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
	okxApiUrl = "https://www.okx.com/api/v5/market/ticker?instId=%s-USDT" // BTC-USDT
)

// Struct to map the JSON response from OKX
type OKXResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		InstType  string `json:"instType"`
		InstId    string `json:"instId"`
		Last      string `json:"last"`
		LastSz    string `json:"lastSz"`
		AskPx     string `json:"askPx"`
		AskSz     string `json:"askSz"`
		BidPx     string `json:"bidPx"`
		BidSz     string `json:"bidSz"`
		Open24h   string `json:"open24h"`
		High24h   string `json:"high24h"`
		Low24h    string `json:"low24h"`
		VolCcy24h string `json:"volCcy24h"`
		Vol24h    string `json:"vol24h"`
		Ts        string `json:"ts"`
		SodUtc0   string `json:"sodUtc0"`
		SodUtc8   string `json:"sodUtc8"`
	} `json:"data"`
}

func GetCoinPriceOkx(coinName string) float64 {
	// URL for the OKX API endpoint
	url := fmt.Sprintf(okxApiUrl, coinName)

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
	var okxResponse OKXResponse
	err = json.Unmarshal(body, &okxResponse)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		os.Exit(1)
	}

	// Check if the response contains data
	if len(okxResponse.Data) == 0 {
		fmt.Println("No data found in the response.")
		os.Exit(1)
	}

	// Get the latest price of BTC
	coinPrice, _ := strconv.ParseFloat(okxResponse.Data[0].Last, 64)
	return coinPrice
}
