package dto

import (
	"cryptotrack/db"
)

type Exchange struct {
	Id   int
	Name string
}

func GetExchangeIdByName(exchangeName string) int {
	database, _ := db.GetSQLiteDBConnection("db.sqlite3")
	defer database.Close()
	result := database.QueryRow("SELECT id FROM exchanges WHERE name = ?", exchangeName)

	var exchangeId int
	result.Scan(&exchangeId)

	return exchangeId
}

func GetAllExchanges() []Exchange {
	database, _ := db.GetSQLiteDBConnection("db.sqlite3")
	defer database.Close()

	result, _ := database.Query("SELECT id, name FROM exchanges")

	exchanges := []Exchange{}
	for result.Next() {
		var e Exchange
		result.Scan(&e.Id, &e.Name)
		exchanges = append(exchanges, e)
	}

	return exchanges
}
