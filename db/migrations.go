package db

func InitMigrations() {

	database, _ := GetSQLiteDBConnection("./db.sqlite3")
	defer database.Close()

	statementExchanges, _ := database.Prepare(`
	CREATE TABLE IF NOT EXISTS exchanges (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	name TEXT
	);`)

	statementExchanges.Exec()

	statementCoins, _ := database.Prepare(`
	CREATE TABLE IF NOT EXISTS coins (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		price REAL DEFAULT 0.0 NOT NULL,
		exchangeId INTEGER,
		UNIQUE(name, exchangeId)
		FOREIGN KEY (exchangeId) REFERENCES exchanges(id)
	)`)

	statementCoins.Exec()

	statementInvestments, _ := database.Prepare(`
	CREATE TABLE IF NOT EXISTS investments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		coinId INTEGER,
		date date,
		investmentInUSD REAL,
		purchasePrice REAL,
		FOREIGN KEY (coinId) REFERENCES coins(id)
	)`)

	statementInvestments.Exec()

	statementStatistics, _ := database.Prepare(`
	CREATE TABLE IF NOT EXISTS statistics (
		metric TEXT,
		value TEXT
	)`)

	statementStatistics.Exec()

	initExchangesData, _ := database.Prepare(`
	INSERT OR IGNORE INTO exchanges (id, name) VALUES (1, "Binance"), (2, "OKX"), (3, "BingX"), (4, "Bybit"), (5, "Gateio")`)

	initExchangesData.Exec()

	addActiveToInvestment, _ := database.Prepare(`
	ALTER TABLE investments ADD COLUMN active BOOLEAN DEFAULT 1;`)

	addActiveToInvestment.Exec()

	addActiveToCoins, _ := database.Prepare(`
	ALTER TABLE coins ADD COLUMN active BOOLEAN DEFAULT 1;`)

	addActiveToCoins.Exec()

	addUpdateDateToCoins, _ := database.Prepare(`
	ALTER TABLE coins ADD COLUMN updateDate date;`)

	addUpdateDateToCoins.Exec()

	statementBalances, _ := database.Prepare(`
	CREATE TABLE IF NOT EXISTS balances (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		currency TEXT,
		balance REAL,
		balanceUSDT REAL,
		exchangeId INTEGER,
		date date,
		timing TEXT
	)`)

	statementBalances.Exec()

	addMetricsBalance, _ := database.Prepare(`INSERT INTO statistics (metric) VALUES ('LastUpdateBalance');`)

	addMetricsBalance.Exec()

}
