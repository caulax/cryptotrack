package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	bybit "github.com/wuhewuhe/bybit.go.api"
)

const (
	bybitApiUrl = "https://api.bybit.com/v5/market/tickers?category=spot&symbol=%sUSDT" // BTCUSDT
)

type BybitCredentials struct {
	ApiKey    string `toml:"apiKey"`
	SecretKey string `toml:"secretKey"`
}

func LoadBybitCredentials(filePath string) (*BybitCredentials, error) {
	var bybitCredentials BybitCredentials
	_, err := toml.DecodeFile(filePath, &struct {
		ByBit *BybitCredentials `toml:"bybit"`
	}{ByBit: &bybitCredentials})
	if err != nil {
		return nil, err
	}
	return &bybitCredentials, nil
}

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
	} else {
		fmt.Println("No ticker data available.")
	}

	var returnValue float64
	// Get the latest price of coin
	returnValue, _ = strconv.ParseFloat(bid1Price, 64)

	return returnValue
}

type AccountBalanceResultBybit struct {
	Currency    string
	Balance     float64
	BalanceUSDT float64
}

func updateBalanceBybit(balanceRes *[]AccountBalanceResultBybit, currency string, newBalance float64, newBalanceUSDT float64) {
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
	*balanceRes = append(*balanceRes, AccountBalanceResultBybit{
		Currency:    currency,
		Balance:     newBalance,
		BalanceUSDT: newBalanceUSDT,
	})
}

func GetWalletBalanceBybit() []AccountBalanceResultBybit {
	config, _ := LoadBybitCredentials("config.toml")

	client := bybit.NewBybitHttpClient(config.ApiKey, config.SecretKey, bybit.WithBaseURL(bybit.MAINNET))

	paramsUnifiedAccount := map[string]interface{}{"accountType": "UNIFIED"}
	accountResultUnified, _ := client.NewUtaBybitServiceWithParams(paramsUnifiedAccount).GetAccountWallet(context.Background())

	var balanceRes []AccountBalanceResultBybit

	if resultSlice, ok := accountResultUnified.Result.(map[string]interface{}); ok {
		if list, ok := resultSlice["list"].([]interface{}); ok {
			for _, account := range list {
				if accountMap, ok := account.(map[string]interface{}); ok {
					if coins, ok := accountMap["coin"].([]interface{}); ok {
						for _, coin := range coins {
							if coinMap, ok := coin.(map[string]interface{}); ok {
								usdValueStr, _ := coinMap["usdValue"].(string)
								usdValueFloat, _ := strconv.ParseFloat(usdValueStr, 64)
								if usdValueFloat > 0.1 {

									walletBalanceStr, _ := coinMap["walletBalance"].(string)
									walletBalanceFloat, _ := strconv.ParseFloat(walletBalanceStr, 64)

									updateBalanceBybit(&balanceRes, coinMap["coin"].(string), walletBalanceFloat, usdValueFloat)
								}
							}
						}
					}
				}
			}
		}
	}

	paramsFundAccount := map[string]interface{}{"accountType": "FUND"}
	accountResultFund, _ := client.NewUtaBybitServiceWithParams(paramsFundAccount).GetAllCoinsBalance(context.Background())

	if resultSlice, ok := accountResultFund.Result.(map[string]interface{}); ok {
		if list, ok := resultSlice["balance"].([]interface{}); ok {
			for _, account := range list {

				if accountMap, ok := account.(map[string]interface{}); ok {

					currency := accountMap["coin"].(string)
					walletBalanceStr, _ := accountMap["walletBalance"].(string)
					walletBalanceFloat, _ := strconv.ParseFloat(walletBalanceStr, 64)

					if currency == "USDT" {
						updateBalanceBybit(&balanceRes, currency, walletBalanceFloat, walletBalanceFloat)
					} else {
						coinPrice := GetCoinPriceBybit(currency)
						balanceUSDT := walletBalanceFloat * coinPrice
						if balanceUSDT > 0.1 {
							updateBalanceBybit(&balanceRes, currency, walletBalanceFloat, balanceUSDT)
						}
					}
				}
			}
		}
	}

	return balanceRes
}
