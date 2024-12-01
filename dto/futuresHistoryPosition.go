package dto

import (
	"cryptotrack/db"
	"database/sql"
	"fmt"
	"time"
)

type FuturesPositionsHistory struct {
	ExchangeName      string
	OpenPositionTime  string
	ClosePositionTime string
	OpenPrice         float64
	ClosePrice        float64
	Leverage          int
	PositionMode      string
	PositionSide      string
	Profit            float64
	CurrencyIn        string
	CurrencyFrom      string
	Fee               float64
	Volume            float64
	TimeInPosition    string
}

func formatDuration(duration time.Duration) string {
	totalSeconds := int(duration.Seconds())
	days := totalSeconds / (24 * 3600)
	hours := (totalSeconds % (24 * 3600)) / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	}
	return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
}

func getDiffDateFrom2Date(dateStart time.Time, dateEnd time.Time) time.Duration {
	diff := dateEnd.Sub(dateStart)
	return diff
}

func GetFuturesHistoryPositionByExchangeId(exchangeId int) ([]FuturesPositionsHistory, float64) {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	var overallProfit float64

	rows, _ := database.Query(`
	SELECT 
		e.Name as exchangeName,
		fHP.openPositionTime as openPositionTime,
		fHP.closePositionTime as closePositionTime,
		fHP.openPrice as openPrice,
		fHP.closePrice as closePrice,
		fHP.leverage as leverage,
		fHP.positionMode as positionMode,
		fHP.positionSide as positionSide,
		fHP.profit as profit,
		fHP.currencyIn as currencyIn,
		fHP.currencyFrom as currencyFrom,
		fHP.fee as fee,
		fHP.volume as volume,
		fHP.timeInPosition as timeInPosition
	FROM futuresHistoryPosition as fHP
	JOIN exchanges AS e ON e.id = fHP.exchangeId
	WHERE fHP.exchangeId = ?
	ORDER BY fHP.openPositionTime DESC
	`, exchangeId)

	futuresPositionsHistory := []FuturesPositionsHistory{}
	for rows.Next() {
		var (
			openPositionTime  int64
			closePositionTime int64
			timeInPosition    int64
			profit            float64
			fPH               FuturesPositionsHistory
		)
		// Scanning fields
		err := rows.Scan(
			&fPH.ExchangeName,
			&openPositionTime,
			&closePositionTime,
			&fPH.OpenPrice,
			&fPH.ClosePrice,
			&fPH.Leverage,
			&fPH.PositionMode,
			&fPH.PositionSide,
			&profit,
			&fPH.CurrencyIn,
			&fPH.CurrencyFrom,
			&fPH.Fee,
			&fPH.Volume,
			&timeInPosition,
		)
		if err != nil {
			fmt.Println(err)
		}

		// Convert Unix times to readable formats
		fPH.OpenPositionTime = time.Unix(openPositionTime/1000, 0).Format("2006-01-02 15:04:05")
		fPH.ClosePositionTime = time.Unix(closePositionTime/1000, 0).Format("2006-01-02 15:04:05")

		fPH.Profit = profit

		fPH.TimeInPosition = formatDuration(
			getDiffDateFrom2Date(
				time.Unix(openPositionTime/1000, 0),
				time.Unix(closePositionTime/1000, 0)))

		overallProfit = overallProfit + profit

		futuresPositionsHistory = append(futuresPositionsHistory, fPH)
	}

	return futuresPositionsHistory, overallProfit
}

func CheckIfTimeExsistInFuturesHistoryPosition(exchangeId int, openPositionTime int64, closePositionTime int64) bool {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	query := fmt.Sprintf(`
	SELECT 
    	1
	FROM futuresHistoryPosition 
	WHERE 
		%d in (SELECT openPositionTime FROM futuresHistoryPosition WHERE exchangeId=%d) AND
		%d in (SELECT closePositionTime FROM futuresHistoryPosition WHERE exchangeId=%d)
	LIMIT 1
	`, openPositionTime, exchangeId, closePositionTime, exchangeId)

	var exists int
	err := database.QueryRow(query).Scan(&exists)
	if err == sql.ErrNoRows {
		return true
	} else if err != nil {
		return true
	}
	return false
}

func CreateNewFuturesHistoryPosition(
	exchangeId int,
	openPositionTime int64,
	closePositionTime int64,
	closePrice float64,
	openPrice float64,
	leverage float64,
	positionMode string,
	positionSide string,
	profit float64,
	currencyIn string,
	currencyFrom string,
	fee float64,
	volume float64,
	timeInPosition int64,
) {
	database, _ := db.GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	statement, _ := database.Prepare(`
	INSERT INTO futuresHistoryPosition (
		exchangeId, 
		openPositionTime, 
		closePositionTime, 
		closePrice, 
		openPrice, 
		leverage, 
		positionMode, 
		positionSide, 
		profit, 
		currencyIn, 
		currencyFrom, 
		fee, 
		volume,
		timeInPosition
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	statement.Exec(
		exchangeId,
		openPositionTime,
		closePositionTime,
		closePrice,
		openPrice,
		leverage,
		positionMode,
		positionSide,
		profit,
		currencyIn,
		currencyFrom,
		fee,
		volume,
		timeInPosition,
	)
}
