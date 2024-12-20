package update

import (
	"cryptotrack/dto"
	"cryptotrack/exchange"
	"fmt"
	"time"
)

func UpdatePrices() {

	coins := dto.GetAllActiveCoins()

	fmt.Println("[INFO] Started update coin prices from exchanges")

	for _, coin := range coins {

		var coinPrice float64

		switch coin.ExchangeId {
		case 1:
			coinPrice = exchange.GetCoinPriceBinance(coin.Name)
		case 2:
			coinPrice = exchange.GetCoinPriceOkx(coin.Name)
		case 3:
			coinPrice = exchange.GetCoinPriceBingx(coin.Name)
		case 4:
			coinPrice = exchange.GetCoinPriceBybit(coin.Name)
		case 5:
			coinPrice = exchange.GetCoinPriceGateio(coin.Name)
		}

		fmt.Println("[INFO] ", coin.Id, coin.Name, coin.ExchangeId, coinPrice)

		coinUpdateTime := time.Now().Local()

		dto.UpdatePriceAndDateOfCoinByIdAndExchangeId(
			dto.GetCoinIdByNameAndExchangeId(
				coin.Name,
				coin.ExchangeId),
			coinPrice,
			coinUpdateTime,
			coin.ExchangeId)
	}

	current_time := time.Now().Local()
	dto.UpdateMetricLastUpdate(current_time.Format("2006-01-02 15:04:05"))

}
