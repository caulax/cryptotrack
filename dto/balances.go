package dto

import (
	"cryptotrack/db"
	"time"
)

type Balance struct {
	Id          int
	CoinName    string
	Balance     float64
	BalanceUSDT float64
	ExchangeId  int
	Date        time.Time
	Timing      string
}

type BalanceOverall struct {
	BalanceUSDT  float64
	ExchangeName string
	Date         time.Time
}

func CreateNewBalance(coinName string, balance float64, balanceUSDT float64, exchangeId int, date time.Time, timing string) {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	statement, _ := database.Prepare("INSERT INTO balances (currency, balance, balanceUSDT, exchangeId, date, timing) VALUES (?, ?, ?, ?, ?, ?)")
	statement.Exec(coinName, balance, balanceUSDT, exchangeId, date, timing)
}

func GetLatestOverallBalanceByTiming(timing string) []BalanceOverall {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	result, _ := database.Query(`
	SELECT 
		balanceUSDT,
		exchangeName,
		date
	FROM (
		SELECT
			SUM(b.balanceUSDT) as balanceUSDT,
			e.name as exchangeName,
			b.date as date,
			b.timing as timing
		FROM balances as b
		JOIN exchanges as e ON e.id = b.exchangeId
		WHERE timing = ?
		GROUP BY date, exchangeId
	) as allBalances
	WHERE date = (SELECT MAX(date) FROM balances WHERE timing = ?)
	`, timing, timing)

	balances := []BalanceOverall{}
	for result.Next() {
		var bal BalanceOverall
		result.Scan(&bal.BalanceUSDT, &bal.ExchangeName, &bal.Date)
		balances = append(balances, bal)
	}

	return balances
}

// func GetAllBalances() []Coin {
// 	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
// 	defer database.Close()

// 	result, _ := database.Query("SELECT id, name, price, exchangeId, active, updateDate FROM coins")

// 	coins := []Coin{}
// 	for result.Next() {
// 		var c Coin
// 		result.Scan(&c.Id, &c.Name, &c.Price, &c.ExchangeId, &c.Active, &c.UpdateDate)
// 		coins = append(coins, c)
// 	}

// 	return coins
// }

// func GetAllActiveBalances() []Coin {
// 	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
// 	defer database.Close()

// 	result, _ := database.Query("SELECT id, name, price, exchangeId, active, updateDate FROM coins WHERE active = 1")

// 	coins := []Coin{}
// 	for result.Next() {
// 		var c Coin
// 		result.Scan(&c.Id, &c.Name, &c.Price, &c.ExchangeId, &c.Active, &c.UpdateDate)
// 		coins = append(coins, c)
// 	}

// 	return coins
// }
