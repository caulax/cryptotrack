package dto

import (
	"cryptotrack/db"
)

type Coin struct {
	Id         int
	Name       string
	Price      float64
	ExchangeId int
}

func CreateNewCoin(coinName string, exchangeId int) {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	statement, _ := database.Prepare("INSERT INTO coins (name, exchangeId) VALUES (?, ?)")
	statement.Exec(coinName, exchangeId)
}

func DeleteCoinById(id int) {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	statement, _ := database.Prepare("DELETE FROM coins WHERE id = ?")
	statement.Exec(id)
}

func UpdatePriceOfCoinByIdAndExchangeId(id int, price float64, exchangeId int) {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	statement, _ := database.Prepare("UPDATE coins SET price = ? WHERE id = ? AND exchangeId = ?")
	statement.Exec(price, id, exchangeId)
}

func GetCoinById(id int) *Coin {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	var c Coin

	result := database.QueryRow("SELECT id, name, price, exchangeId FROM coins WHERE id = ?", id)
	result.Scan(&c.Id, &c.Name, &c.Price, &c.ExchangeId)

	return &c
}

func GetAllCoins() []Coin {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	result, _ := database.Query("SELECT id, name, price, exchangeId FROM coins")

	coins := []Coin{}
	for result.Next() {
		var c Coin
		result.Scan(&c.Id, &c.Name, &c.Price, &c.ExchangeId)
		coins = append(coins, c)
	}

	return coins
}

func GetCoinIdByNameAndExchangeId(name string, exchangeId int) int {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	result := database.QueryRow("SELECT id FROM coins WHERE name = ? AND exchangeId = ?", name, exchangeId)

	var id int
	result.Scan(&id)

	return id
}
