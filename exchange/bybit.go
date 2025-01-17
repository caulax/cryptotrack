package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

func LoadBybitCredentials(filePath, credentialsBlock string) (*BybitCredentials, error) {
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

	// Unmarshal the JSON into the BybitCredentials struct
	var bybitCredentials BybitCredentials
	err = json.Unmarshal(jsonData, &bybitCredentials)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling block data into bybitCredentials: %w", err)
	}

	return &bybitCredentials, nil
}

// func LoadBybitCredentials(filePath string) (*BybitCredentials, error) {
// 	var bybitCredentials BybitCredentials
// 	_, err := toml.DecodeFile(filePath, &struct {
// 		ByBit *BybitCredentials `toml:"bybit"`
// 	}{ByBit: &bybitCredentials})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &bybitCredentials, nil
// }

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
		fmt.Println(coinName, "No ticker data available.")
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

func GetWalletBalanceBybit(accountCredentialsBlock string) []AccountBalanceResultBybit {
	config, _ := LoadBybitCredentials("config.toml", accountCredentialsBlock)

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

func GetCurrentTime() int64 {
	now := time.Now()
	unixNano := now.UnixNano()
	timeStamp := unixNano / int64(time.Millisecond)
	return timeStamp
}

// Trade execution response struct
type TradeExecution struct {
	OrderID  string  `json:"orderId"`
	Fee      float64 `json:"execFee,string"`
	Symbol   string  `json:"symbol"`
	Side     string  `json:"side"`
	ExecType string  `json:"execType"`
}

// API Response structure for the `/v5/execution/list` endpoint
type ExecutionResponse struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		Executions []TradeExecution `json:"list"`
	} `json:"result"`
}

func GetTradeFeeByOrderID(orderID, accountCredentialsBlock string) float64 {

	config, _ := LoadBybitCredentials("config.toml", accountCredentialsBlock)

	client := bybit.NewBybitHttpClient(
		config.ApiKey,
		config.SecretKey,
		bybit.WithBaseURL(bybit.MAINNET),
	)

	// Define the parameters for the API request
	params := map[string]interface{}{
		"orderId":  orderID, // Specify the order ID
		"category": "linear",
	}

	response, _ := client.NewUtaBybitServiceWithParams(params).GetTradeHistory(context.Background())

	jsonString, _ := json.Marshal(response)

	var executionResponse ExecutionResponse
	json.Unmarshal(jsonString, &executionResponse)

	// Return the list of executions
	return executionResponse.Result.Executions[0].Fee
}

// Input struct to match the "list" items in JSON
type ClosePnlEntry struct {
	CreatedTime   string `json:"createdTime"`
	UpdatedTime   string `json:"updatedTime"`
	AvgExitPrice  string `json:"avgExitPrice"`
	AvgEntryPrice string `json:"avgEntryPrice"`
	Leverage      string `json:"leverage"`
	OrderType     string `json:"orderType"`
	OrderId       string `json:"orderId"`
	Side          string `json:"side"`
	ClosedPnl     string `json:"closedPnl"`
	ClosedSize    string `json:"closedSize"`
	OrderPrice    string `json:"orderPrice"`
	Symbol        string `json:"symbol"`
}

// Output struct with transformed fields
type PositionsHistoryBybit struct {
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

func TransformClosePnlEntries(entries []ClosePnlEntry, accountCredentialsBlock string) []PositionsHistoryBybit {
	var transformedEntries []PositionsHistoryBybit

	for _, entry := range entries {
		// Parse required fields
		createdTime, _ := strconv.ParseInt(entry.CreatedTime, 10, 64)
		updatedTime, _ := strconv.ParseInt(entry.UpdatedTime, 10, 64)
		closePrice, _ := strconv.ParseFloat(entry.AvgExitPrice, 64)
		openPrice, _ := strconv.ParseFloat(entry.AvgEntryPrice, 64)
		leverage, _ := strconv.ParseFloat(entry.Leverage, 64)
		closedPnl, _ := strconv.ParseFloat(entry.ClosedPnl, 64)
		closedSize, _ := strconv.ParseFloat(entry.ClosedSize, 64)
		orderPrice, _ := strconv.ParseFloat(entry.OrderPrice, 64)

		// Split symbol into currencyFrom and currencyIn
		parts := strings.Split(entry.Symbol, "USDT")
		currencyIn := parts[0]
		currencyFrom := "USDT"

		// Map "side" to "positionSide"
		positionSide := "short"
		if entry.Side == "Buy" {
			positionSide = "long"
		}

		// Calculate volume and time in position
		volume := (closedSize * orderPrice) / leverage
		timeInPosition := updatedTime - createdTime

		// Append transformed entry
		transformedEntries = append(transformedEntries, PositionsHistoryBybit{
			OpenPositionTime:  createdTime,
			ClosePositionTime: updatedTime,
			ClosePrice:        closePrice,
			OpenPrice:         openPrice,
			Leverage:          leverage,
			PositionMode:      entry.OrderType,
			PositionSide:      positionSide,
			Profit:            closedPnl,
			CurrencyIn:        currencyIn,
			CurrencyFrom:      currencyFrom,
			Fee:               GetTradeFeeByOrderID(entry.OrderId, accountCredentialsBlock),
			Volume:            volume,
			TimeInPosition:    timeInPosition,
		})
	}

	return transformedEntries
}

func GetWalletPositionsHistoryBybit(accountCredentialsBlock string) []PositionsHistoryBybit {
	config, _ := LoadBybitCredentials("config.toml", accountCredentialsBlock)

	client := bybit.NewBybitHttpClient(
		config.ApiKey,
		config.SecretKey,
		bybit.WithBaseURL(bybit.MAINNET),
	)

	paramsUnifiedAccount := map[string]interface{}{
		"accountType": "UNIFIED",
		"category":    "linear",
		"limit":       100,
		// "startTime": 1729415216000,
	}
	closePnl, _ := client.NewUtaBybitServiceWithParams(paramsUnifiedAccount).GetClosePnl(context.Background())

	jsonString, _ := json.Marshal(closePnl)

	var input struct {
		Result struct {
			List []ClosePnlEntry `json:"list"`
		} `json:"result"`
	}

	err := json.Unmarshal([]byte(jsonString), &input)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
	}

	positionsHistoryBybit := TransformClosePnlEntries(input.Result.List, accountCredentialsBlock)

	// for _, entry := range positionsHistoryBybit {
	// 	fmt.Println(entry.OpenPositionTime, entry.ClosePositionTime, entry.ClosePrice, entry.OpenPrice, entry.Leverage, entry.PositionMode, entry.PositionSide, entry.Profit, entry.CurrencyIn, entry.CurrencyFrom, entry.Fee, entry.Volume, entry.TimeInPosition)
	// }

	return positionsHistoryBybit

}

// todo: rewrite bybit api to http lib, without github.com/wuhewuhe/bybit.go.api
// func GetEarnWalletBalanceBybit(accountCredentialsBlock string) {
// 	config, _ := LoadBybitCredentials("config.toml", accountCredentialsBlock)

// 	req, _ := http.NewRequest("GET", "https://api.bybit.com"+"/v5/broker/earnings-info", nil)

// 	recvWindow := "5000"

// 	timeStamp := GetCurrentTime()

// 	signatureBase := []byte(strconv.FormatInt(timeStamp, 10) + config.ApiKey + recvWindow)
// 	hmac256 := hmac.New(sha256.New, []byte(config.SecretKey))
// 	hmac256.Write(signatureBase)
// 	signature := hex.EncodeToString(hmac256.Sum(nil))

// 	req.Header.Add("X-BAPI-SIGN", signature)
// 	req.Header.Add("X-BAPI-SIGN-TYPE", "2")
// 	req.Header.Add("X-BAPI-API-KEY", config.ApiKey)
// 	req.Header.Add("X-BAPI-TIMESTAMP", strconv.FormatInt(timeStamp, 10))
// 	req.Header.Add("X-BAPI-RECV-WINDOW", recvWindow)
// 	req.Header.Add("User-Agent", fmt.Sprintf("%s/%s", "bybit.api.go", "1.0.4"))

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error making request:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Read and print the response
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response body:", err)
// 		return
// 	}

// 	fmt.Println("Response Status:", resp.Status)
// 	fmt.Println("Response Body:", string(body))
// }
