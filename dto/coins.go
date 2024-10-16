package dto

import (
	"cryptotrack/db"
	"time"
)

type Coin struct {
	Id         int
	Name       string
	Price      float64
	ExchangeId int
	Active     bool
	UpdateDate time.Time
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

func UpdatePriceAndDateOfCoinByIdAndExchangeId(id int, price float64, updateDate time.Time, exchangeId int) {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	statement, _ := database.Prepare("UPDATE coins SET price = ?, updateDate = ? WHERE id = ? AND exchangeId = ?")
	statement.Exec(price, updateDate, id, exchangeId)
}

func GetCoinById(id int) *Coin {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	var c Coin

	result := database.QueryRow("SELECT id, name, price, exchangeId, active, updateDate FROM coins WHERE id = ?", id)
	result.Scan(&c.Id, &c.Name, &c.Price, &c.ExchangeId, &c.Active, &c.UpdateDate)

	return &c
}

func GetAllCoins() []Coin {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	result, _ := database.Query("SELECT id, name, price, exchangeId, active, updateDate FROM coins")

	coins := []Coin{}
	for result.Next() {
		var c Coin
		result.Scan(&c.Id, &c.Name, &c.Price, &c.ExchangeId, &c.Active, &c.UpdateDate)
		coins = append(coins, c)
	}

	return coins
}

func GetAllActiveCoins() []Coin {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	result, _ := database.Query("SELECT id, name, price, exchangeId, active, updateDate FROM coins WHERE active = 1")

	coins := []Coin{}
	for result.Next() {
		var c Coin
		result.Scan(&c.Id, &c.Name, &c.Price, &c.ExchangeId, &c.Active, &c.UpdateDate)
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

func ActivateCoinById(id int) {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	statement, _ := database.Prepare("UPDATE coins SET active = 1 WHERE id = ?")
	statement.Exec(id)
}

func DeactivateCoinById(id int) {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	statement, _ := database.Prepare("UPDATE coins SET active = 0 WHERE id = ?")
	statement.Exec(id)
}
