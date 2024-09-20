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

	initExchangesData, _ := database.Prepare(`
	INSERT OR IGNORE INTO exchanges (id, name) VALUES (1, "Binance"), (2, "OKX"), (3, "BingX")`)

	initExchangesData.Exec()

	addActiveToInvestment, _ := database.Prepare(`
	ALTER TABLE investments ADD COLUMN active BOOLEAN DEFAULT 1;`)

	addActiveToInvestment.Exec()

}
