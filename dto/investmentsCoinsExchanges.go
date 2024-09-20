package dto

import (
	"cryptotrack/db"
	"time"
)

type InvestmentsCoinsExchanges struct {
	Id              int
	Date            time.Time
	InvestmentInUSD float64
	PurchasePrice   float64
	Active          bool
	CoinName        string
	CurrentPrice    float64
	ExchangeName    string
}

func GetInvestmentsCoinsExchanges() []InvestmentsCoinsExchanges {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	result, _ := database.Query(`
		SELECT 
			i.id as id,
			i.date as date,
			i.investmentInUSD as investmentInUSD,
			i.purchasePrice as purchasePrice,
			i.active as active,
			c.name as coinName,
			c.price as currentPrice,
			e.name as exchangeName
		FROM 
			investments AS i
		JOIN 
			coins AS c ON i.coinId = c.id
		JOIN 
			exchanges AS e ON c.exchangeId = e.id
		ORDER BY i.date
	`)

	investments := []InvestmentsCoinsExchanges{}
	for result.Next() {
		var i InvestmentsCoinsExchanges
		result.Scan(&i.Id, &i.Date, &i.InvestmentInUSD, &i.PurchasePrice, &i.Active, &i.CoinName, &i.CurrentPrice, &i.ExchangeName)
		investments = append(investments, i)
	}

	return investments
}
