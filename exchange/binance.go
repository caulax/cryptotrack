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
	binanceApiUrl = "https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT" // BTCUSDT
)

// Struct to map the JSON response from Binance
type BinanceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func GetCoinPriceBinance(coinName string) float64 {
	// URL for the Binance API endpoint
	url := fmt.Sprintf(binanceApiUrl, coinName)

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
	var binanceResponse BinanceResponse
	err = json.Unmarshal(body, &binanceResponse)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		os.Exit(1)
	}

	// Check if the response contains data
	if binanceResponse == (BinanceResponse{}) {
		fmt.Println(coinName, "No ticker data available.")
	}

	// Get the latest price of BTC
	coinPrice, _ := strconv.ParseFloat(binanceResponse.Price, 64)
	return coinPrice
}
