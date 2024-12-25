package service

import (
	"cryptotrack/dto"
)

type OverallInformation struct {
	AmountInvestmentInUSD float64
	AmountIncome          float64
	AmountOverallIncome   float64
	ExchangeData          map[string][]ExchangeDetails
}

type ExchangeDetails struct {
	Id              int
	Date            string
	InvestmentInUSD float64
	PurchasePrice   float64
	Active          bool
	CoinName        string
	CurrentPrice    float64
	Profit          int
	Income          float64
	OverallIncome   float64
}

func CountProfitIncomeOverallIncome() map[string][]ExchangeDetails {
	exchangeDetails := make(map[string][]ExchangeDetails)
	investments := dto.GetInvestmentsCoinsExchanges()

	for _, v := range investments {
		if v.Active {
			detail := ExchangeDetails{
				Id:              v.Id,
				Date:            v.Date.Format("2006-01-02"),
				InvestmentInUSD: v.InvestmentInUSD,
				PurchasePrice:   v.PurchasePrice,
				Active:          v.Active,
				CoinName:        v.CoinName,
				CurrentPrice:    v.CurrentPrice,
				Profit:          int((v.CurrentPrice - v.PurchasePrice) / v.PurchasePrice * 100),
				Income:          float64(v.InvestmentInUSD * (float64(int((v.CurrentPrice-v.PurchasePrice)/v.PurchasePrice*100)) / 100)),
				OverallIncome:   v.InvestmentInUSD + float64(v.InvestmentInUSD*(float64(int((v.CurrentPrice-v.PurchasePrice)/v.PurchasePrice*100))/100)),
			}
			exchangeDetails[v.ExchangeName] = append(exchangeDetails[v.ExchangeName], detail)
		}
	}
	return exchangeDetails
}

func CountProfitIncomeOverallIncomeArchive() map[string][]ExchangeDetails {
	exchangeDetails := make(map[string][]ExchangeDetails)
	investments := dto.GetInvestmentsCoinsExchanges()

	for _, v := range investments {
		if !v.Active {
			detail := ExchangeDetails{
				Id:              v.Id,
				Date:            v.Date.Format("2006-01-02"),
				InvestmentInUSD: v.InvestmentInUSD,
				PurchasePrice:   v.PurchasePrice,
				Active:          v.Active,
				CoinName:        v.CoinName,
				CurrentPrice:    v.CurrentPrice,
				Profit:          int((v.CurrentPrice - v.PurchasePrice) / v.PurchasePrice * 100),
				Income:          float64(v.InvestmentInUSD * (float64(int((v.CurrentPrice-v.PurchasePrice)/v.PurchasePrice*100)) / 100)),
				OverallIncome:   v.InvestmentInUSD + float64(v.InvestmentInUSD*(float64(int((v.CurrentPrice-v.PurchasePrice)/v.PurchasePrice*100))/100)),
			}
			exchangeDetails[v.ExchangeName] = append(exchangeDetails[v.ExchangeName], detail)
		}
	}
	return exchangeDetails
}

func GetOverallInformation() OverallInformation {
	exchangeData := CountProfitIncomeOverallIncome()
	overallInformation := OverallInformation{ExchangeData: exchangeData}

	for _, details := range exchangeData {
		for _, v := range details {
			overallInformation.AmountInvestmentInUSD += v.InvestmentInUSD
			overallInformation.AmountIncome += v.Income
			overallInformation.AmountOverallIncome += v.OverallIncome
		}
	}
	return overallInformation
}

func GetArchiveInformation() OverallInformation {
	exchangeData := CountProfitIncomeOverallIncomeArchive()
	overallInformation := OverallInformation{ExchangeData: exchangeData}

	for _, details := range exchangeData {
		for _, v := range details {
			overallInformation.AmountInvestmentInUSD += v.InvestmentInUSD
			overallInformation.AmountIncome += v.Income
			overallInformation.AmountOverallIncome += v.OverallIncome
		}
	}
	return overallInformation
}
