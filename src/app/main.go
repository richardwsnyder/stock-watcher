package main

import (
	"database/sql"
	"encoding/json"
	"finnhub/src/database"
	d "finnhub/src/database"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Stock struct {
	Symbol string
	Name   sql.NullString
	Price  float64
}

func WaitForCtrlC() {
	var end_waiter sync.WaitGroup
	end_waiter.Add(1)
	var signal_channel chan os.Signal
	signal_channel = make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)
	go func() {
		<-signal_channel
		end_waiter.Done()
	}()
	end_waiter.Wait()
}

func getQuote(symbol string, token string) map[string]interface{} {
	client := http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=brpoutnrh5rbpquqb4s0", symbol), nil)
	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	return result
}

func (s *Stock) updateStock(db *sql.DB, newValue float64) {
	fmt.Printf("Updating stock %v to price %v\n", s.Symbol, s.Price)
	sqlStatement := `
	UPDATE stocks
	SET price = $1
	WHERE symbol = $2
	`

	res, err := db.Exec(sqlStatement, s.Price, s.Symbol)
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println(count)
}

func (s *Stock) insertStock(db *sql.DB) {
	quote := getQuote(s.Symbol, d.Token)

	s.Price = quote["c"].(float64)

	sqlStatement := `
	INSERT INTO stocks (Symbol, Price)
	VALUES ($1, $2)
	`

	_, err := db.Exec(sqlStatement, s.Symbol, s.Price)
	if err != nil {
		panic(err)
	}
}

func watchStock(stock *Stock, db *sql.DB) {
	fmt.Println("Beginning to watch stock", stock.Symbol)
	for {
		quote := getQuote(stock.Symbol, d.Token)

		oldPrice := stock.Price
		stock.Price = quote["c"].(float64)

		var newPrice float64
		sqlStatement := `
		SELECT price FROM stocks WHERE symbol = $1
		`

		row := db.QueryRow(sqlStatement, stock.Symbol)
		switch err := row.Scan(&newPrice); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned for symbol", stock.Symbol)
		case nil:
			if newPrice != oldPrice {
				stock.updateStock(db, newPrice)
			}
		default:
			panic(err)
		}

		time.Sleep(1 * time.Minute)
	}
}

func main() {
	db, err := database.Connect(database.Host, database.Port, database.User, database.Password, database.Dbname)

	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT symbol, name, price FROM stocks")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var symbol string
		var name sql.NullString
		var price float64

		err = rows.Scan(&symbol, &name, &price)
		if err != nil {
			panic(err)
		}
		s := Stock{
			Symbol: symbol,
			Name:   name,
			Price:  price,
		}

		go watchStock(&s, db)
	}

	WaitForCtrlC()
}
