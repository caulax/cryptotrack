package update

import (
	"cryptotrack/dto"
	"cryptotrack/exchange"
	"fmt"
	"time"
)

func UpdateFuturesHistoryPostion() {

	fmt.Println("[INFO] Starting update futures history position for OKX")
	okxID := dto.GetExchangeIdByName("OKX")
	futuresHistoryPositionsOkx := exchange.GetWalletPositionsHistoryOkx()
	for _, futuresHistoryPositionOkx := range futuresHistoryPositionsOkx {
		if dto.CheckIfTimeExsistInFuturesHistoryPosition(okxID, futuresHistoryPositionOkx.OpenPositionTime, futuresHistoryPositionOkx.ClosePositionTime) {
			dto.CreateNewFuturesHistoryPosition(
				okxID,
				futuresHistoryPositionOkx.OpenPositionTime,
				futuresHistoryPositionOkx.ClosePositionTime,
				futuresHistoryPositionOkx.ClosePrice,
				futuresHistoryPositionOkx.OpenPrice,
				futuresHistoryPositionOkx.Leverage,
				futuresHistoryPositionOkx.PositionMode,
				futuresHistoryPositionOkx.PositionSide,
				futuresHistoryPositionOkx.Profit,
				futuresHistoryPositionOkx.CurrencyIn,
				futuresHistoryPositionOkx.CurrencyFrom,
				futuresHistoryPositionOkx.Fee,
				futuresHistoryPositionOkx.Volume,
				futuresHistoryPositionOkx.TimeInPosition,
			)
		}
	}
	fmt.Println("[INFO] Updated futures history position for OKX")

	fmt.Println("[INFO] Starting update futures history position for Bybit")
	bybitID := dto.GetExchangeIdByName("Bybit")
	futuresHistoryPositionsBybit := exchange.GetWalletPositionsHistoryBybit()
	for _, futuresHistoryPositionBybit := range futuresHistoryPositionsBybit {
		if dto.CheckIfTimeExsistInFuturesHistoryPosition(bybitID, futuresHistoryPositionBybit.OpenPositionTime, futuresHistoryPositionBybit.ClosePositionTime) {
			dto.CreateNewFuturesHistoryPosition(
				bybitID,
				futuresHistoryPositionBybit.OpenPositionTime,
				futuresHistoryPositionBybit.ClosePositionTime,
				futuresHistoryPositionBybit.ClosePrice,
				futuresHistoryPositionBybit.OpenPrice,
				futuresHistoryPositionBybit.Leverage,
				futuresHistoryPositionBybit.PositionMode,
				futuresHistoryPositionBybit.PositionSide,
				futuresHistoryPositionBybit.Profit,
				futuresHistoryPositionBybit.CurrencyIn,
				futuresHistoryPositionBybit.CurrencyFrom,
				futuresHistoryPositionBybit.Fee,
				futuresHistoryPositionBybit.Volume,
				futuresHistoryPositionBybit.TimeInPosition,
			)
		}
	}
	fmt.Println("[INFO] Updated futures history position for Bybit")

	current_time := time.Now().Local()
	dto.UpdateMetricLastUpdateFuturesHistoryPosition(current_time.Format("2006-01-02 15:04:05"))

}
