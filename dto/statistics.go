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

func UpdateMetricLastUpdate(date string) {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	statement, _ := database.Prepare("UPDATE statistics SET value = ? WHERE metric = 'LastUpdate'")
	statement.Exec(date)
}
