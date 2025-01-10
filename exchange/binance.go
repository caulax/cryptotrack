package exchange

import (
	"crypto/hmac"
	"crypto/sha256"
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

const (
	binanceAPIBaseURL = "https://api.binance.com"
	accountEndpoint   = "/api/v3/account"
)

type AccountResponseBinance struct {
	Balances []struct {
		Asset  string `json:"asset"`
		Free   string `json:"free"`
		Locked string `json:"locked"`
	} `json:"balances"`
}

type AccountBalanceResultBinance struct {
	Currency    string
	Balance     float64
	BalanceUSDT float64
}

func mapAssetName(asset string) string {
	if strings.HasPrefix(asset, "LD") {
		return strings.TrimPrefix(asset, "LD")
	}
	return asset
}

func updateBalanceBinance(balanceRes *[]AccountBalanceResultBinance, currency string, newBalance float64, newBalanceUSDT float64) {
	// Loop through the slice to find the currency
	for i, bal := range *balanceRes {
		if bal.Currency == currency {
			// Update the Balance and BalanceUSDT
			(*balanceRes)[i].Balance = (*balanceRes)[i].Balance + newBalance
			(*balanceRes)[i].BalanceUSDT = (*balanceRes)[i].BalanceUSDT + newBalanceUSDT
			return
		}
	}

	// If currency not found, append a new entry
	*balanceRes = append(*balanceRes, AccountBalanceResultBinance{
		Currency:    currency,
		Balance:     newBalance,
		BalanceUSDT: newBalanceUSDT,
	})
}

func signQuery(query, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(query))
	return hex.EncodeToString(h.Sum(nil))
}

type BinanceCredentials struct {
	ApiKey    string `toml:"apiKey"`
	SecretKey string `toml:"secretKey"`
}

func LoadBinanceCredentials(filePath, credentialsBlock string) (*BinanceCredentials, error) {
	var tomlData map[string]interface{}

	// Decode the TOML file into a generic map
	_, err := toml.DecodeFile(filePath, &tomlData)
	if err != nil {
		return nil, err
	}

	// Retrieve the desired block
	blockData, ok := tomlData[credentialsBlock]
	if !ok {
		return nil, fmt.Errorf("block %s not found in TOML file", credentialsBlock)
	}

	// Marshal the block data back into JSON (for compatibility with structs)
	jsonData, err := json.Marshal(blockData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling block data: %w", err)
	}

	// Unmarshal the JSON into the binanceCredentials struct
	var binanceCredentials BinanceCredentials
	err = json.Unmarshal(jsonData, &binanceCredentials)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling block data into binance credentials: %w", err)
	}

	return &binanceCredentials, nil
}

func GetWalletBalanceBinance(accountCredentialsBlock string) []AccountBalanceResultBinance {
	config, _ := LoadBinanceCredentials("config.toml", accountCredentialsBlock)

	timestamp := time.Now().UnixMilli()
	query := fmt.Sprintf("timestamp=%d", timestamp)
	signature := signQuery(query, config.SecretKey)
	query = fmt.Sprintf("%s&signature=%s", query, signature)

	req, err := http.NewRequest("GET", binanceAPIBaseURL+accountEndpoint+"?"+query, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		// return
	}

	req.Header.Set("X-MBX-APIKEY", config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		// return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Received status code %d\n", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response: %s\n", string(body))
		// return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		// return
	}
	var accountResponse AccountResponseBinance
	err = json.Unmarshal(body, &accountResponse)
	if err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		// return
	}

	var results []AccountBalanceResultBinance
	for _, balance := range accountResponse.Balances {
		freeBalance, err := strconv.ParseFloat(balance.Free, 64)
		if err != nil {
			fmt.Printf("Error converting balance: %v\n", err)
			continue
		}
		if freeBalance > 0 {
			currency := mapAssetName(balance.Asset)

			if currency == "USDT" {
				updateBalanceBinance(
					&results,
					currency,
					freeBalance,
					freeBalance,
				)
			} else {
				currencyBalanceUSDT := freeBalance * GetCoinPriceBinance(currency)

				if currencyBalanceUSDT > 0.1 {
					updateBalanceBinance(
						&results,
						currency,
						freeBalance,
						currencyBalanceUSDT,
					)
				}
			}
		}
	}

	return results

}
