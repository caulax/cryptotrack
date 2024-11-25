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
