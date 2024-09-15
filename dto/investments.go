package dto

import (
	"cryptotrack/db"
	"time"
)

type Investment struct {
	id              int
	coinId          int
	date            time.Time
	investmentInUSD float64
	purchasePrice   float64
}

func CreateNewInvestment(coinId int, date time.Time, investmentInUSD float64, purchasePrice float64) {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	statement, _ := database.Prepare("INSERT INTO investments (coinId, date, investmentInUSD, purchasePrice) VALUES (?, ?, ?, ?)")
	statement.Exec(coinId, date, investmentInUSD, purchasePrice)
}

func DeleteInvestmentById(id int) {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	statement, _ := database.Prepare("DELETE FROM investments WHERE id = ?")
	statement.Exec(id)
}

func GetAllInvestment() []Investment {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	result, _ := database.Query("SELECT id, coinId, date, investmentInUSD, purchasePrice FROM investments")

	investments := []Investment{}
	for result.Next() {
		var i Investment
		result.Scan(&i.id, &i.coinId, &i.date, &i.investmentInUSD, &i.purchasePrice)
		investments = append(investments, i)
	}

	return investments
}
