package service

import (
	"cryptotrack/dto"
	"time"
)

type BalanceOverallTable struct {
	ExchangeName          string
	BalanceUSDT           float64
	BalanceUSDTHourly     float64
	BalanceUSDTHourlyDiff float64
	BalanceUSDTDaily      float64
	BalanceUSDTDailyDiff  float64
}

type BalanceByCoin struct {
	CoinName              string
	Balance               float64
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

func GetAllBalancesWithDiffByExchangeAndCoin() map[string][]BalanceByCoin {
	balanceCoinByExchange := make(map[string][]BalanceByCoin)

	balance := dto.GetLatestBalanceByTiming("minute")
	balanceHourly := dto.GetLatestBalanceByTiming("hourly")
	balanceDaily := dto.GetLatestBalanceByTiming("daily")

	hourlyMap := make(map[string]float64)
	dailyMap := make(map[string]float64)

	for _, bh := range balanceHourly {
		hourlyMap[bh.ExchangeName+bh.CoinName] = bh.BalanceUSDT
	}
	for _, bd := range balanceDaily {
		dailyMap[bd.ExchangeName+bd.CoinName] = bd.BalanceUSDT
	}

	for _, b := range balance {
		key := b.ExchangeName + b.CoinName
		hourlyDiff := 0.0
		if hourlyMap[key] != 0 {
			hourlyDiff = b.BalanceUSDT - hourlyMap[key]
		}
		dailyDiff := 0.0
		if dailyMap[key] != 0 {
			dailyDiff = b.BalanceUSDT - dailyMap[key]
		}

		detail := BalanceByCoin{
			CoinName:              b.CoinName,
			Balance:               b.Balance,
			BalanceUSDT:           b.BalanceUSDT,
			BalanceUSDTHourly:     hourlyMap[key],
			BalanceUSDTHourlyDiff: hourlyDiff,
			BalanceUSDTDaily:      dailyMap[key],
			BalanceUSDTDailyDiff:  dailyDiff,
		}
		balanceCoinByExchange[b.ExchangeName] = append(balanceCoinByExchange[b.ExchangeName], detail)
	}

	return balanceCoinByExchange
}

func CleanUpBalances() {
	dto.DeleteBalanceByDate(time.Now().AddDate(0, 0, -30))
}
