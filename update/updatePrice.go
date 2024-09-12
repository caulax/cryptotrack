package update

import (
	"cryptotrack/dto"
	"cryptotrack/exchange"
	"fmt"
)

func UpdatePrices() {

	coins := dto.GetAllCoins()

	for _, coin := range coins {

		var coinPrice float64

		switch coin.ExchangeId {
		case 1:
			coinPrice = exchange.GetCoinPriceBinance(coin.Name)
		case 2:
			coinPrice = exchange.GetCoinPriceOkx(coin.Name)
		case 3:
			coinPrice = exchange.GetCoinPriceBingx(coin.Name)
		}

		fmt.Println(coin.Name, coin.ExchangeId, coinPrice)

		dto.UpdatePriceOfCoinByIdAndExchangeId(
			dto.GetCoinIdByNameAndExchangeId(
				coin.Name,
				coin.ExchangeId),
			coinPrice,
			coin.ExchangeId)
	}

}