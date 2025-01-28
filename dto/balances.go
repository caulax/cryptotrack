package dto

import (
	"cryptotrack/db"
	"fmt"
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

type BalanceByCoin struct {
	CoinName     string
	Balance      float64
	BalanceUSDT  float64
	ExchangeName string
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
		SUM(maxDate.balanceUSDT) as balanceUSDT,
		e.name as ExchangeName,
		maxDate.date as Date
	FROM (
		SELECT
			b.balanceUSDT,
			b.exchangeId,
			b.date
		FROM balances as b
		WHERE date = (SELECT MAX(date) FROM balances WHERE timing = ?)
		) as maxDate
	JOIN exchanges as e ON e.id = maxDate.exchangeId
	GROUP BY maxDate.date, maxDate.exchangeId
	ORDER BY balanceUSDT DESC
	`, timing)

	balances := []BalanceOverall{}
	for result.Next() {
		var bal BalanceOverall
		result.Scan(&bal.BalanceUSDT, &bal.ExchangeName, &bal.Date)
		balances = append(balances, bal)
	}

	return balances
}

func GetLatestBalanceByTiming(timing string) []BalanceByCoin {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	result, _ := database.Query(`
	SELECT 
		maxDate.currency as CoinName,
		maxDate.balance as Balance,
		maxDate.balanceUSDT as BalanceUSDT,
		e.name as ExchangeName
	FROM (
		SELECT
			b.currency,
			b.balance,
			b.balanceUSDT,
			b.exchangeId
		FROM balances AS b
		WHERE date = (SELECT MAX(date) FROM balances WHERE timing = ?)
		) as maxDate
	JOIN exchanges AS e ON e.id = maxDate.exchangeId
	ORDER BY BalanceUSDT DESC
	`, timing)

	balances := []BalanceByCoin{}
	for result.Next() {
		var bal BalanceByCoin
		result.Scan(&bal.CoinName, &bal.Balance, &bal.BalanceUSDT, &bal.ExchangeName)
		balances = append(balances, bal)
	}

	return balances

}

func DeleteBalanceByDate(date time.Time) {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	// Prepare the DELETE query
	query := "DELETE FROM balances WHERE timing = 'minute' AND date < ?"
	statement, err := database.Prepare(query)
	if err != nil {
		fmt.Printf("Failed to prepare statement: %v", err)
	}
	defer statement.Close()

	// Execute the DELETE query
	result, err := statement.Exec(date.Format("2006-01-02"))
	if err != nil {
		fmt.Printf("Failed to execute query: %v", err)
	}

	// Get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Failed to get rows affected: %v", err)
	}

	fmt.Printf("Deleted %d rows older than 30 days.\n", rowsAffected)

}

func GetAllTimeDailyBalance() []BalanceOverall {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	result, _ := database.Query(`
		SELECT
			e.name as exchangeName,
			date,
			SUM(balanceUSDT)
		FROM balances
		JOIN exchanges AS e ON exchangeId = e.id
		WHERE timing = 'daily'
		GROUP BY exchangeId, date
		ORDER BY date ASC;
	`)

	balances := []BalanceOverall{}
	for result.Next() {
		var bal BalanceOverall
		result.Scan(&bal.ExchangeName, &bal.Date, &bal.BalanceUSDT)
		balances = append(balances, bal)
	}

	return balances
}
