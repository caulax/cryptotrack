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

	tableData := service.GetOverallInformation()
	fmt.Println("[INFO] Overall Info:", tableData)
	err = t.Execute(w, tableData)
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

func NewServer() *server {
	s := &server{
		subscriberMessageBuffer: 10,
		subscribers:             make(map[*subscriber]struct{}),
	}

	fileServer := http.FileServer(http.Dir("./static"))
	s.mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	s.mux.HandleFunc("/activate/{id}", handlerActivateInvestment)
	s.mux.HandleFunc("/deactivate/{id}", handlerDectivateInvestment)
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
	case "update":
		update.UpdatePrices()
	case "migrations":
		db.InitMigrations()
	default:
		fmt.Println("Invalid mode. Use 'server' or 'update'.")
		os.Exit(1)
	}
}
