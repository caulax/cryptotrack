package dto

import (
	"cryptotrack/db"
	"time"
)

type CoinsExchanges struct {
	Id         int
	Name       string
	Price      float64
	Exchange   string
	Active     bool
	UpdateDate time.Time
}

func GetAllCoinsExchanges() []CoinsExchanges {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	result, _ := database.Query(`
		SELECT 
			c.id, 
			c.name, 
			c.price, 
			e.name, 
			c.active, 
			c.updateDate 
		FROM 
			coins AS c 
		JOIN 
			exchanges AS e ON c.exchangeId = e.id;`)

	coinsExchanges := []CoinsExchanges{}
	for result.Next() {
		var c CoinsExchanges
		result.Scan(&c.Id, &c.Name, &c.Price, &c.Exchange, &c.Active, &c.UpdateDate)
		coinsExchanges = append(coinsExchanges, c)
	}

	return coinsExchanges
}
