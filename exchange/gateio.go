package exchange

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	gateioApiUrl = "https://api.gateio.ws/api/v4/spot/tickers?currency_pair=%s_USDT" // BTC_USDT

	baseURLGateio = "https://api.gateio.ws"
	prefixGateio  = "/api/v4"
	methodGateio  = "GET"
	urlPathGateio = "/spot/accounts"
)

type GateioCredentials struct {
	ApiKey    string `toml:"apiKey"`
	ApiSecret string `toml:"apiSecret"`
}

func LoadGateioCredentials(filePath string) (*GateioCredentials, error) {
	var gateioCredentials GateioCredentials
	_, err := toml.DecodeFile(filePath, &struct {
		Gateio *GateioCredentials `toml:"gateio"`
	}{Gateio: &gateioCredentials})
	if err != nil {
		return nil, err
	}
	return &gateioCredentials, nil
}

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

type AccountBalanceGateio struct {
	Currency  string `json:"currency"`
	Available string `json:"available"`
	Locked    string `json:"locked"`
	UpdateID  int    `json:"update_id"`
}

type AccountBalanceResultGateio struct {
	Currency    string
	Balance     float64
	BalanceUSDT float64
}

func GetWalletBalanceGateio() []AccountBalanceResultGateio {
	config, _ := LoadGateioCredentials("config.toml")
	queryParam := ""
	bodyParam := ""

	// Generate timestamp
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	// Generate body hash
	bodyHash := sha512Hash(bodyParam)

	// Generate sign string
	signString := fmt.Sprintf("%s\n%s%s\n%s\n%s\n%s", methodGateio, prefixGateio, urlPathGateio, queryParam, bodyHash, timestamp)

	// Generate HMAC-SHA512 signature
	signature := hmacSha512(signString, config.ApiSecret)

	// Construct full URL
	fullURL := baseURLGateio + prefixGateio + urlPathGateio

	// Make HTTP request
	client := &http.Client{}
	req, err := http.NewRequest(methodGateio, fullURL, bytes.NewBuffer([]byte(bodyParam)))
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	// Add headers
	req.Header.Add("Timestamp", timestamp)
	req.Header.Add("KEY", config.ApiKey)
	req.Header.Add("SIGN", signature)
	req.Header.Add("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}

	// Map response to struct
	var accounts []AccountBalanceGateio
	err = json.Unmarshal(body, &accounts)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	var accountBalanceResult []AccountBalanceResultGateio
	for _, account := range accounts {

		var accountBalance AccountBalanceResultGateio

		if account.Currency == "USDT" {
			accountBalance.Currency = account.Currency
			accountBalance.Balance, _ = strconv.ParseFloat(strings.TrimSpace(account.Available), 64)
			accountBalance.BalanceUSDT, _ = strconv.ParseFloat(strings.TrimSpace(account.Available), 64)
			accountBalanceResult = append(accountBalanceResult, accountBalance)
		} else {
			accountBalance.Currency = account.Currency
			accountBalance.Balance, _ = strconv.ParseFloat(strings.TrimSpace(account.Available), 64)
			accountBalance.BalanceUSDT = accountBalance.Balance * GetCoinPriceGateio(account.Currency)
			if accountBalance.BalanceUSDT > 0.1 {
				accountBalanceResult = append(accountBalanceResult, accountBalance)
			}
		}
	}

	return accountBalanceResult
}

// Helper function to calculate SHA512 hash
func sha512Hash(data string) string {
	hash := sha512.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// Helper function to calculate HMAC-SHA512 signature
func hmacSha512(data, secret string) string {
	h := hmac.New(sha512.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
