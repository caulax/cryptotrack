package service

import (
	"cryptotrack/dto"
	"time"
)

func FormatTimestamp(ms int64) string {
	sec := ms / 1000 // Convert milliseconds to seconds
	return time.Unix(sec, 0).Format("2006-01-02 15:04:05")
}

func GetDiffDate() (time.Duration, bool) {

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	lastUpdateTime := dto.GetMetricLastUpdate().Value

	parsedCurrentTime, _ := time.Parse("2006-01-02 15:04:05", currentTime)
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", lastUpdateTime)

	diff := parsedCurrentTime.Sub(parsedTime)

	timeAlert := false

	if diff > 10*time.Minute {
		timeAlert = true
	}

	return diff, timeAlert

}

func GetDiffDateBalance() (time.Duration, bool) {

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	lastUpdateTime := dto.GetMetricLastUpdateBalance().Value

	parsedCurrentTime, _ := time.Parse("2006-01-02 15:04:05", currentTime)
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", lastUpdateTime)

	diff := parsedCurrentTime.Sub(parsedTime)

	timeAlert := false

	if diff > 10*time.Minute {
		timeAlert = true
	}

	return diff, timeAlert

}

func GetDiffDateFuturesHistoryPosition() (time.Duration, bool) {

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	lastUpdateTime := dto.GetMetricLastUpdateFuturesHistoryPosition().Value

	parsedCurrentTime, _ := time.Parse("2006-01-02 15:04:05", currentTime)
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", lastUpdateTime)

	diff := parsedCurrentTime.Sub(parsedTime)

	timeAlert := false

	if diff > 10*time.Minute {
		timeAlert = true
	}

	return diff, timeAlert

}

func GetDiffDateFromDate(date time.Time) time.Duration {

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	parsedCurrentTime, _ := time.Parse("2006-01-02 15:04:05", currentTime)

	diff := parsedCurrentTime.Sub(date)

	return diff

}
