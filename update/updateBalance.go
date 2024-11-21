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

	balanceOkx := exchange.GetWalletBalanceOkx()
	for _, balance := range balanceOkx {
		dto.CreateNewBalance(balance.Currency, balance.Balance, balance.BalanceUSDT, 2, updateTime, timing)
	}
	fmt.Println("[INFO] OKX Balance: ", balanceOkx)

	balanceBybit := exchange.GetWalletBalanceBybit()
	for _, balance := range balanceBybit {
		dto.CreateNewBalance(balance.Currency, balance.Balance, balance.BalanceUSDT, 4, updateTime, timing)
	}
	fmt.Println("[INFO] ByBit Balance: ", balanceBybit)

	balanceGateio := exchange.GetWalletBalanceGateio()
	for _, balance := range balanceGateio {
		dto.CreateNewBalance(balance.Currency, balance.Balance, balance.BalanceUSDT, 5, updateTime, timing)
	}
	fmt.Println("[INFO] GateIO Balance: ", balanceGateio)

	current_time := time.Now().Local()
	dto.UpdateMetricLastUpdateBalance(current_time.Format("2006-01-02 15:04:05"))
	
	fmt.Println("[INFO] Balances from exchanges saved to db")
}
