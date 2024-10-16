package service

import (
	"cryptotrack/dto"
	"time"
)

func GetDiffDate() time.Duration {

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	lastUpdateTime := dto.GetMetricLastUpdate().Value

	parsedCurrentTime, _ := time.Parse("2006-01-02 15:04:05", currentTime)
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", lastUpdateTime)

	diff := parsedCurrentTime.Sub(parsedTime)

	return diff

}

func GetDiffDateFromDate(date time.Time) time.Duration {

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	parsedCurrentTime, _ := time.Parse("2006-01-02 15:04:05", currentTime)

	diff := parsedCurrentTime.Sub(date)

	return diff

}
