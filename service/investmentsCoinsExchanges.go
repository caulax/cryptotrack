package service

import (
	"cryptotrack/dto"
)

type OverallInformation struct {
	AmountInvestmentInUSD            float64
	AmountIncome                     float64
	AmountOverallIncome              float64
	ServiceInvestmentsCoinsExchanges []ServiceInvestmentsCoinsExchanges
}

type ServiceInvestmentsCoinsExchanges struct {
	Date            string
	InvestmentInUSD float64
	PurchasePrice   float64
	CoinName        string
	CurrentPrice    float64
	ExchangeName    string
	Profit          int
	Income          float64
	OverallIncome   float64
}

func CountProfitIncomeOverallIncome() []ServiceInvestmentsCoinsExchanges {

	serviceInvestments := []ServiceInvestmentsCoinsExchanges{}
	investments := dto.GetInvestmentsCoinsExchanges()

	for _, v := range investments {
		var sICE ServiceInvestmentsCoinsExchanges

		sICE.Date = v.Date.Format("2006-01-02")
		sICE.InvestmentInUSD = v.InvestmentInUSD
		sICE.PurchasePrice = v.PurchasePrice
		sICE.CoinName = v.CoinName
		sICE.CurrentPrice = v.CurrentPrice
		sICE.ExchangeName = v.ExchangeName
		sICE.Profit = int((v.CurrentPrice - v.PurchasePrice) / v.PurchasePrice * 100)
		sICE.Income = float64(v.InvestmentInUSD * (float64(sICE.Profit) / 100))
		sICE.OverallIncome = v.InvestmentInUSD + sICE.Income

		serviceInvestments = append(serviceInvestments, sICE)

	}
	return serviceInvestments
}

func GetOverallInformation() OverallInformation {

	var overallInformation OverallInformation

	profitIncomeOverallIncome := CountProfitIncomeOverallIncome()

	for _, v := range profitIncomeOverallIncome {
		overallInformation.AmountInvestmentInUSD += v.InvestmentInUSD
		overallInformation.AmountIncome += v.Income
		overallInformation.AmountOverallIncome += v.OverallIncome
	}

	overallInformation.ServiceInvestmentsCoinsExchanges = profitIncomeOverallIncome

	return overallInformation
}
