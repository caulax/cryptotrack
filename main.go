package main

import (
	"context"
	"cryptotrack/db"
	"cryptotrack/dto"
	"cryptotrack/service"
	"cryptotrack/update"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type subscriber struct {
	msgs chan []byte
}

type server struct {
	subscriberMessageBuffer int
	mux                     http.ServeMux
	subscriberMutex         sync.Mutex
	subscribers             map[*subscriber]struct{}
}

type Variables struct {
	TableData         service.OverallInformation
	TableDataBalance  []service.BalanceOverallTable
	LastUpdate        time.Duration
	TimeAlert         bool
	LastUpdateBalance time.Duration
	TimeAlertBalance  bool
}

type VariablesCoin struct {
	TableData  []service.CoinsExchanges
	LastUpdate time.Duration
	TimeAlert  bool
}

type VariablesFutures struct {
	TableData     []dto.FuturesPositionsHistory
	LastUpdate    time.Duration
	TimeAlert     bool
	OverallProfit float64
}

func (s *server) subscriberHandler(w http.ResponseWriter, r *http.Request) {
	err := s.subscribe(r.Context(), w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *server) addSubscriber(subscriber *subscriber) {
	s.subscriberMutex.Lock()
	s.subscribers[subscriber] = struct{}{}
	s.subscriberMutex.Unlock()
	fmt.Println("[INFO] Added subscriber", subscriber)
}

func (s *server) subscribe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var c *websocket.Conn
	subscriber := &subscriber{
		msgs: make(chan []byte, s.subscriberMessageBuffer),
	}
	s.addSubscriber(subscriber)
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}
	defer c.CloseNow()
	ctx = c.CloseRead(ctx)

	for {
		select {
		case msg := <-subscriber.msgs:
			ctx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			err := c.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func handlerMain(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("htmx/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	diff, timeAlert := service.GetDiffDate()
	diffBalance, timeAlertBalance := service.GetDiffDateBalance()

	balanceOverall := service.GetAllBalancesWithDiff()

	tableData := service.GetOverallInformation()
	fmt.Println("[INFO] Overall Info:", tableData)

	v := Variables{
		TableData:         tableData,
		TableDataBalance:  balanceOverall,
		LastUpdate:        diff,
		TimeAlert:         timeAlert,
		LastUpdateBalance: diffBalance,
		TimeAlertBalance:  timeAlertBalance,
	}

	err = t.Execute(w, &v)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
	}
}

func handlerAddNewCrypto(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		coinName := r.FormValue("coin-name")
		purchasePrice := r.FormValue("purchase-price")
		datePurchase := r.FormValue("date-purchase")
		investmentInUsd := r.FormValue("investment-in-usd")
		exchangeId := r.FormValue("exchange")

		exchangeIdInt, _ := strconv.Atoi(exchangeId)

		coinExchangeId := dto.GetCoinIdByNameAndExchangeId(coinName, exchangeIdInt)

		if coinExchangeId == 0 {
			dto.CreateNewCoin(coinName, exchangeIdInt)
			coinExchangeId = dto.GetCoinIdByNameAndExchangeId(coinName, exchangeIdInt)
		}

		datePurchaseTime, _ := time.Parse("2006-01-02", datePurchase)
		investmentInUsdFloat, _ := strconv.ParseFloat(investmentInUsd, 64)
		purchasePriceFloat, _ := strconv.ParseFloat(purchasePrice, 64)

		dto.CreateNewInvestment(coinExchangeId, datePurchaseTime, investmentInUsdFloat, purchasePriceFloat)

		http.Redirect(w, r, "/", http.StatusSeeOther)

	case http.MethodGet:
		t, err := template.ParseFiles("htmx/add.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		exchangeList := dto.GetAllExchanges()

		fmt.Println(exchangeList)

		err = t.Execute(w, exchangeList)
		if err != nil {
			fmt.Println(err)
			w.Write([]byte(err.Error()))
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func handlerUpdate(w http.ResponseWriter, r *http.Request) {
	update.UpdatePrices()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handlerCoin(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("htmx/coins.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	diff, timeAlert := service.GetDiffDate()

	tableData := service.GetAllCoinsExchangesWithDiffTime()
	fmt.Println("[INFO] Overall Info:", tableData)

	v := VariablesCoin{
		TableData:  tableData,
		LastUpdate: diff,
		TimeAlert:  timeAlert,
	}

	err = t.Execute(w, &v)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
	}
}

func handlerArchive(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("htmx/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	diff, timeAlert := service.GetDiffDate()

	tableData := service.GetArchiveInformation()
	fmt.Println("[INFO] Overall Info:", tableData)

	v := Variables{
		TableData:  tableData,
		LastUpdate: diff,
		TimeAlert:  timeAlert,
	}

	err = t.Execute(w, &v)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
	}
}

func handlerActivateCoin(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/coin/activate/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	dto.ActivateCoinById(id)
	http.Redirect(w, r, "/coins", http.StatusSeeOther)
}

func handlerDectivateCoin(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/coin/deactivate/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	dto.DeactivateCoinById(id)
	http.Redirect(w, r, "/coins", http.StatusSeeOther)
}

func handlerActivateInvestment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/activate/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	dto.ActivateInvestementById(id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handlerDectivateInvestment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/deactivate/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	dto.DeactivateInvestementById(id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handlerFuturesHistoryPositionByExchange(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/futures/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		id = 2
	}

	t, err := template.ParseFiles("htmx/futures.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	diff, timeAlert := service.GetDiffDateFuturesHistoryPosition()

	tableData, overallProfit := dto.GetFuturesHistoryPositionByExchangeId(id)
	fmt.Println("[INFO] Futures Info:", tableData)

	v := VariablesFutures{
		TableData:     tableData,
		LastUpdate:    diff,
		TimeAlert:     timeAlert,
		OverallProfit: overallProfit,
	}

	err = t.Execute(w, &v)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
	}
}

func NewServer() *server {
	s := &server{
		subscriberMessageBuffer: 10,
		subscribers:             make(map[*subscriber]struct{}),
	}

	fileServer := http.FileServer(http.Dir("./static"))
	s.mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	s.mux.HandleFunc("/activate/{id}", handlerActivateInvestment)
	s.mux.HandleFunc("/deactivate/{id}", handlerDectivateInvestment)
	s.mux.HandleFunc("/coin/activate/{id}", handlerActivateCoin)
	s.mux.HandleFunc("/coin/deactivate/{id}", handlerDectivateCoin)
	s.mux.HandleFunc("/futures/{exchange}", handlerFuturesHistoryPositionByExchange)
	s.mux.HandleFunc("/archive", handlerArchive)
	s.mux.HandleFunc("/coins", handlerCoin)
	s.mux.HandleFunc("/update", handlerUpdate)
	s.mux.HandleFunc("/add", handlerAddNewCrypto)
	s.mux.HandleFunc("/", handlerMain)
	s.mux.HandleFunc("/ws", s.subscriberHandler)
	return s
}

func (s *server) broadcast(msg []byte) {
	s.subscriberMutex.Lock()
	for subscriber := range s.subscribers {
		subscriber.msgs <- msg
	}
	s.subscriberMutex.Unlock()
}

func main() {

	mode := flag.String("mode", "server", "Mode to run: 'server' or 'update'")
	flag.Parse()

	switch *mode {
	case "server":
		srv := NewServer()
		err := http.ListenAndServe(":8080", &srv.mux)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "updatePrices":
		update.UpdatePrices()
	case "updateBalance":
		update.UpdateBalance("minute")
	case "updateBalanceHourly":
		update.UpdateBalance("hourly")
	case "updateBalanceDaily":
		update.UpdateBalance("daily")
	case "updateFuturesHistoryPostion":
		update.UpdateFuturesHistoryPostion()
	case "cleanUpBalances":
		service.CleanUpBalances()
	case "migrations":
		db.InitMigrations()
	// case "test":
	// exchange.GetWalletPositionsHistoryBybit()
	// exchange.GetWalletPositionsHistoryOkx("v-okx")
	// // fmt.Println("OKX: ", exchange.GetWalletBalanceOkx())
	// // fmt.Println("GATEIO: ", exchange.GetWalletBalanceGateio())
	// // fmt.Println("BYBIT: ", exchange.GetWalletBalanceBybit())
	// // fmt.Println(dto.GetLatestOverallBalanceByTiming("minute"))
	// // fmt.Println(service.GetAllBalancesWithDiff())
	default:
		fmt.Println("Invalid mode. Use 'server' or 'update'.")
		os.Exit(1)
	}
}
