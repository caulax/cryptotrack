package update

import (
	"cryptotrack/dto"
	"cryptotrack/exchange"
	"fmt"
	"time"
)

func UpdateBalance(timing string) {

	updateTime := time.Now().Local()

	fmt.Println("[INFO] Started getting balances from exchanges")

	balanceBinance := exchange.GetWalletBalanceBinance("binance")
	for _, balance := range balanceBinance {
		dto.CreateNewBalance(balance.Currency, balance.Balance, balance.BalanceUSDT, 1, updateTime, timing)
	}
	fmt.Println("[INFO] Binance Balance: ", balanceBinance)

	balanceOkx := exchange.GetWalletBalanceOkx("okx")
	for _, balance := range balanceOkx {
		dto.CreateNewBalance(balance.Currency, balance.Balance, balance.BalanceUSDT, 2, updateTime, timing)
	}
	fmt.Println("[INFO] OKX Balance: ", balanceOkx)

	balanceBybit := exchange.GetWalletBalanceBybit("bybit")
	for _, balance := range balanceBybit {
		dto.CreateNewBalance(balance.Currency, balance.Balance, balance.BalanceUSDT, 4, updateTime, timing)
	}
	fmt.Println("[INFO] ByBit Balance: ", balanceBybit)

	balanceGateio := exchange.GetWalletBalanceGateio()
	for _, balance := range balanceGateio {
		dto.CreateNewBalance(balance.Currency, balance.Balance, balance.BalanceUSDT, 5, updateTime, timing)
	}
	fmt.Println("[INFO] GateIO Balance: ", balanceGateio)

	balanceOkxV := exchange.GetWalletBalanceOkx("v-okx")
	for _, balance := range balanceOkxV {
		dto.CreateNewBalance(balance.Currency, balance.Balance, balance.BalanceUSDT, 6, updateTime, timing)
	}
	fmt.Println("[INFO] V-OKX Balance: ", balanceOkxV)

	balanceBybitV := exchange.GetWalletBalanceBybit("v-bybit")
	for _, balance := range balanceBybitV {
		dto.CreateNewBalance(balance.Currency, balance.Balance, balance.BalanceUSDT, 7, updateTime, timing)
	}
	fmt.Println("[INFO] V-Bybit Balance: ", balanceBybitV)

	current_time := time.Now().Local()
	dto.UpdateMetricLastUpdateBalance(current_time.Format("2006-01-02 15:04:05"))

	fmt.Println("[INFO] Balances from exchanges saved to db")
}
