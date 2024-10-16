package service

import (
	"cryptotrack/dto"
	"time"
)

type CoinsExchanges struct {
	Id         int
	Name       string
	Price      float64
	Exchange   string
	Active     bool
	UpdateDate time.Duration
	TimeAlert  bool
}

func GetAllCoinsExchangesWithDiffTime() []CoinsExchanges {

	coinsExchangesDiffTime := []CoinsExchanges{}
	coinsExchanges := dto.GetAllCoinsExchanges()

	for _, v := range coinsExchanges {
		var ce CoinsExchanges

		ce.Id = v.Id
		ce.Name = v.Name
		ce.Price = v.Price
		ce.Exchange = v.Exchange
		ce.Active = v.Active

		ce.UpdateDate = GetDiffDateFromDate(v.UpdateDate)

		ce.TimeAlert = false

		if ce.UpdateDate > 10*time.Minute {
			ce.TimeAlert = true
		}

		coinsExchangesDiffTime = append(coinsExchangesDiffTime, ce)

	}
	return coinsExchangesDiffTime
}
