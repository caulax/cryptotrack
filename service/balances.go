package service

import (
	"cryptotrack/dto"
)

type BalanceOverallTable struct {
	ExchangeName          string
	BalanceUSDT           float64
	BalanceUSDTHourly     float64
	BalanceUSDTHourlyDiff float64
	BalanceUSDTDaily      float64
	BalanceUSDTDailyDiff  float64
}

func GetAllBalancesWithDiff() []BalanceOverallTable {

	balanceUSDT := dto.GetLatestOverallBalanceByTiming("minute")
	balanceUSDTHourly := dto.GetLatestOverallBalanceByTiming("hourly")
	balanceUSDTDaily := dto.GetLatestOverallBalanceByTiming("daily")

	var balanceOverallTableMap []BalanceOverallTable

	for _, balanceU := range balanceUSDT {
		var balanceOverallTable BalanceOverallTable
		for _, balanceHourly := range balanceUSDTHourly {
			for _, balanceDaily := range balanceUSDTDaily {

				if balanceU.ExchangeName == balanceHourly.ExchangeName && balanceU.ExchangeName == balanceDaily.ExchangeName {

					balanceOverallTable.ExchangeName = balanceU.ExchangeName
					balanceOverallTable.BalanceUSDT = balanceU.BalanceUSDT
					balanceOverallTable.BalanceUSDTHourly = balanceHourly.BalanceUSDT
					balanceOverallTable.BalanceUSDTHourlyDiff = balanceU.BalanceUSDT - balanceHourly.BalanceUSDT
					balanceOverallTable.BalanceUSDTDaily = balanceDaily.BalanceUSDT
					balanceOverallTable.BalanceUSDTDailyDiff = balanceU.BalanceUSDT - balanceDaily.BalanceUSDT

					balanceOverallTableMap = append(balanceOverallTableMap, balanceOverallTable)
				}
			}
		}
	}

	return balanceOverallTableMap

}
