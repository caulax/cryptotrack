package dto

import (
	"cryptotrack/db"
)

type Statistics struct {
	Metric string
	Value  string
}

func GetMetricLastUpdate() *Statistics {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	var s Statistics

	result := database.QueryRow("SELECT metric, value FROM statistics WHERE metric = 'LastUpdate'")
	result.Scan(&s.Metric, &s.Value)

	return &s
}

func GetMetricLastUpdateBalance() *Statistics {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	var s Statistics

	result := database.QueryRow("SELECT metric, value FROM statistics WHERE metric = 'LastUpdateBalance'")
	result.Scan(&s.Metric, &s.Value)

	return &s
}

func GetMetricLastUpdateFuturesHistoryPosition() *Statistics {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	var s Statistics

	result := database.QueryRow("SELECT metric, value FROM statistics WHERE metric = 'LastUpdateFuturesHistoryPosition'")
	result.Scan(&s.Metric, &s.Value)

	return &s
}

func UpdateMetricLastUpdate(date string) {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	statement, _ := database.Prepare("UPDATE statistics SET value = ? WHERE metric = 'LastUpdate'")
	statement.Exec(date)
}

func UpdateMetricLastUpdateBalance(date string) {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	statement, _ := database.Prepare("UPDATE statistics SET value = ? WHERE metric = 'LastUpdateBalance'")
	statement.Exec(date)
}

func UpdateMetricLastUpdateFuturesHistoryPosition(date string) {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	statement, _ := database.Prepare("UPDATE statistics SET value = ? WHERE metric = 'LastUpdateFuturesHistoryPosition'")
	statement.Exec(date)
}
