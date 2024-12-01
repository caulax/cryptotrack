package exchange

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
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
	okxApiUrl = "https://www.okx.com/api/v5/market/ticker?instId=%s-USDT" // BTC-USDT

	baseURL     = "https://www.okx.com"
	apiEndpoint = "/api/v5/account/balance"
)

type OkxCredentials struct {
	ApiKey     string `toml:"apiKey"`
	SecretKey  string `toml:"secretKey"`
	Passphrase string `toml:"passphrase"`
}

func LoadOkxCredentials(filePath string) (*OkxCredentials, error) {
	var okxCredentials OkxCredentials
	_, err := toml.DecodeFile(filePath, &struct {
		Okx *OkxCredentials `toml:"okx"`
	}{Okx: &okxCredentials})
	if err != nil {
		return nil, err
	}

	return &okxCredentials, nil
}

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
		fmt.Println(coinName, "No data found in the response.")
	}

	// Get the latest price of BTC
	coinPrice, _ := strconv.ParseFloat(okxResponse.Data[0].Last, 64)
	return coinPrice
}

type AccountBalance struct {
	Code string `json:"code"`
	Data []struct {
		Details []struct {
			Currency string `json:"ccy"`
			Balance  string `json:"cashBal"`
		} `json:"details"`
	} `json:"data"`
}

type AccountBalanceResult struct {
	Currency    string
	Balance     float64
	BalanceUSDT float64
}

type EarnBalance struct {
	Code string `json:"code"`
	Data []struct {
		Currency string `json:"ccy"`
		Amount   string `json:"amt"`
	} `json:"data"`
}

func updateBalanceOkx(balanceRes *[]AccountBalanceResult, currency string, newBalance float64, newBalanceUSDT float64) {
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
	*balanceRes = append(*balanceRes, AccountBalanceResult{
		Currency:    currency,
		Balance:     newBalance,
		BalanceUSDT: newBalanceUSDT,
	})
}

// generateSignature creates a signature for OKX API authentication
func generateSignature(timestamp, method, endpoint, body string) string {
	config, _ := LoadOkxCredentials("config.toml")

	signatureString := timestamp + method + endpoint + body
	h := hmac.New(sha256.New, []byte(config.SecretKey))
	h.Write([]byte(signatureString))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// getRequest performs a GET request to the specified OKX endpoint
func getRequest(endpoint string) (*http.Response, error) {
	config, _ := LoadOkxCredentials("config.toml")

	timestamp := time.Now().UTC().Format(time.RFC3339)
	signature := generateSignature(timestamp, "GET", endpoint, "")

	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Set headers for authentication
	req.Header.Add("OK-ACCESS-KEY", config.ApiKey)
	req.Header.Add("OK-ACCESS-SIGN", signature)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", config.Passphrase)
	req.Header.Add("Content-Type", "application/json")

	return client.Do(req)
}

// getTradingBalance fetches and prints trading account balance
func GetWalletBalanceOkx() []AccountBalanceResult {

	resp, err := getRequest("/api/v5/account/balance")
	if err != nil {
		fmt.Println("Error fetching trading balance:", err)
	}
	defer resp.Body.Close()

	var balance AccountBalance
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		fmt.Println("Error decoding trading balance:", err)
	}

	var balanceRes []AccountBalanceResult
	for _, data := range balance.Data {
		for _, detail := range data.Details {

			balance64, _ := strconv.ParseFloat(strings.TrimSpace(detail.Balance), 64)

			if detail.Currency == "USDT" {
				updateBalanceOkx(&balanceRes, detail.Currency, balance64, balance64)
			} else {
				coinPrice := GetCoinPriceOkx(detail.Currency)
				balanceUSDT := balance64 * coinPrice
				if balanceUSDT > 0.1 {
					updateBalanceOkx(&balanceRes, detail.Currency, balance64, balanceUSDT)
				}
			}
		}
	}

	respEarnFlex, err := getRequest("/api/v5/finance/savings/balance")
	if err != nil {
		fmt.Println("Error fetching earn balance:", err)
	}
	defer respEarnFlex.Body.Close()

	var balanceEarn EarnBalance
	if err := json.NewDecoder(respEarnFlex.Body).Decode(&balanceEarn); err != nil {
		fmt.Println("Error decoding earn balance:", err)
	}

	if len(balanceEarn.Data) > 0 {
		for _, data := range balanceEarn.Data {
			balanceEarn64, _ := strconv.ParseFloat(strings.TrimSpace(data.Amount), 64)
			updateBalanceOkx(&balanceRes, data.Currency, balanceEarn64, balanceEarn64)
		}
	}

	return balanceRes
}

type PositionsHistoryOkx struct {
	OpenPositionTime  int64   `json:"openPositionTime"`
	ClosePositionTime int64   `json:"closePositionTime"`
	ClosePrice        float64 `json:"closePrice"`
	OpenPrice         float64 `json:"openPrice"`
	Leverage          float64 `json:"leverage"`
	PositionMode      string  `json:"positionMode"`
	PositionSide      string  `json:"positionSide"`
	Profit            float64 `json:"profit"`
	CurrencyIn        string  `json:"currencyIn"`
	CurrencyFrom      string  `json:"currencyFrom"`
	Fee               float64 `json:"fee"`
	Volume            float64 `json:"volume"`
	TimeInPosition    int64   `json:"timeInPosition"`
}

func mapToPosition(data map[string]interface{}) PositionsHistoryOkx {
	// Parse raw data fields
	cTime, _ := strconv.ParseInt(data["cTime"].(string), 10, 64)
	uTime, _ := strconv.ParseInt(data["uTime"].(string), 10, 64)
	closeAvgPx, _ := strconv.ParseFloat(data["closeAvgPx"].(string), 64)
	openAvgPx, _ := strconv.ParseFloat(data["openAvgPx"].(string), 64)
	lever, _ := strconv.ParseFloat(data["lever"].(string), 64)
	openMaxPos, _ := strconv.ParseFloat(data["openMaxPos"].(string), 64)
	fee, _ := strconv.ParseFloat(data["fee"].(string), 64)
	fundingFee, _ := strconv.ParseFloat(data["fundingFee"].(string), 64)
	liqPenalty, _ := strconv.ParseFloat(data["liqPenalty"].(string), 64)
	realizedPnl, _ := strconv.ParseFloat(data["realizedPnl"].(string), 64)
	uly := data["uly"].(string)

	// Split `uly` into currencies
	currencies := strings.Split(uly, "-")
	currencyIn := currencies[0]
	currencyFrom := currencies[1]

	// Calculate fee, volume, and time in position
	totalFee := fee + fundingFee + liqPenalty
	volume := (openMaxPos * openAvgPx) / lever
	timeInPosition := uTime - cTime

	// openTimeFormatted := formatTimestamp(cTime)
	// closeTimeFormatted := formatTimestamp(uTime)

	// Map fields to the new struct
	return PositionsHistoryOkx{
		OpenPositionTime:  cTime,
		ClosePositionTime: uTime,
		ClosePrice:        closeAvgPx,
		OpenPrice:         openAvgPx,
		Leverage:          lever,
		PositionMode:      data["mgnMode"].(string),
		PositionSide:      data["posSide"].(string),
		Profit:            realizedPnl,
		CurrencyIn:        currencyIn,
		CurrencyFrom:      currencyFrom,
		Fee:               totalFee,
		Volume:            volume,
		TimeInPosition:    timeInPosition,
	}
}

func GetWalletPositionsHistoryOkx() []PositionsHistoryOkx {

	resp, err := getRequest("/api/v5/account/positions-history")
	if err != nil {
		fmt.Println("Error fetching trading balance:", err)
	}
	defer resp.Body.Close()

	bodyResp, _ := io.ReadAll(resp.Body)

	var data map[string]interface{}
	err = json.Unmarshal(bodyResp, &data)
	if err != nil {
		fmt.Println("Error:", err)
	}

	var positionsHistoryOkx []PositionsHistoryOkx
	dataList := data["data"].([]interface{})
	for _, item := range dataList {
		data := item.(map[string]interface{})
		position := mapToPosition(data)
		positionsHistoryOkx = append(positionsHistoryOkx, position)
	}

	return positionsHistoryOkx
}

// type FundingBalance struct {
// 	Code string `json:"code"`
// 	Data []struct {
// 		Currency string `json:"ccy"`
// 		Balance  string `json:"bal"`
// 	} `json:"data"`
// }

// // getFundingBalance fetches and prints funding account balance
// func GetFundingBalance() {
// 	resp, err := getRequest("/api/v5/asset/balances")
// 	if err != nil {
// 		fmt.Println("Error fetching funding balance:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	var balance FundingBalance
// 	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
// 		fmt.Println("Error decoding funding balance:", err)
// 		return
// 	}

// 	fmt.Println("Funding Account Balance:")
// 	for _, data := range balance.Data {
// 		fmt.Printf("Currency: %s, Balance: %s\n", data.Currency, data.Balance)
// 	}
// }
