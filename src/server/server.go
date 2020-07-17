package server

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	st "github.com/richardwsnyder/stock-watcher/src/stock"
)

type StocksPage struct {
	PageTitle string
	Stocks    []st.Stock
}

func Serve(db *sql.DB) {
	r := mux.NewRouter()
	r.HandleFunc("/", homepage)
	r.HandleFunc("/stock/{symbol}", individualStock)
	r.HandleFunc("/stocks", func(w http.ResponseWriter, r *http.Request) {
		tpl, err := template.ParseFiles("../views/stocks.html")
		if err != nil {
			tpl, err = template.ParseFiles("src/views/stocks.html")
		}
		tmpl := template.Must(tpl, err)
		stocks := allStocks(db)
		stocksPage := StocksPage{
			PageTitle: "All Stocks",
			Stocks:    stocks,
		}
		tmpl.Execute(w, stocksPage)
	})
	fmt.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, you've requested the home page at: %s\n", r.URL.Path)
}

func allStocks(db *sql.DB) []st.Stock {
	rows, err := db.Query("SELECT symbol, name, price, pricetarget FROM stocks ORDER BY price DESC")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	stocks := []st.Stock{}

	for rows.Next() {
		var symbol string
		var name sql.NullString
		var price float64
		var priceTarget sql.NullFloat64

		err = rows.Scan(&symbol, &name, &price, &priceTarget)
		if err != nil {
			panic(err)
		}
		s := st.Stock{
			Symbol:      symbol,
			Name:        name,
			Price:       price,
			PriceTarget: priceTarget,
		}

		stocks = append(stocks, s)
	}

	return stocks
}

func individualStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "This is the stock symbol: %s", vars["symbol"])
}
