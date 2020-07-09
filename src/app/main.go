package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/richardwsnyder/stock-watcher/src/database"
	em "github.com/richardwsnyder/stock-watcher/src/email"
	en "github.com/richardwsnyder/stock-watcher/src/env"
)

var Token = en.GoDotEnvVariable("TOKEN")

type Stock struct {
	Symbol      string
	Name        sql.NullString
	Price       float64
	PriceTarget sql.NullFloat64
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

func getQuote(symbol string) map[string]interface{} {
	client := http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", symbol, Token), nil)
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

	_, err := db.Exec(sqlStatement, s.Price, s.Symbol)
	if err != nil {
		panic(err)
	}
}

func (s *Stock) insertStock(db *sql.DB) {
	quote := getQuote(s.Symbol)

	s.Price = quote["c"].(float64)

	sqlStatement := `
	INSERT INTO stocks (symbol, price, pricetarget, name)
	VALUES ($1, $2, $3, $4)
	`

	_, err := db.Exec(sqlStatement, s.Symbol, s.Price, s.PriceTarget.Float64, s.Name)
	if err != nil {
		panic(err)
	}
}

func updatePriceTarget(symbol string, newPriceTaret float64, db *sql.DB) {
	fmt.Printf("Updating stock %v to new price target %v\n", symbol, newPriceTaret)
	sqlStatement := `
	UPDATE stocks
	SET pricetarget = $1
	WHERE symbol = $2
	`

	_, err := db.Exec(sqlStatement, newPriceTaret, symbol)
	if err != nil {
		panic(err)
	}
}

func removeStock(symbol string, db *sql.DB) {
	fmt.Printf("Removing stock %v\n", symbol)
	sqlStatement := `
	DELETE FROM stocks
	WHERE symbol = $1
	`

	_, err := db.Exec(sqlStatement, symbol)
	if err != nil {
		panic(err)
	}
}

func watchStock(stock *Stock, db *sql.DB) {
	fmt.Println("Beginning to watch stock", stock.Symbol)
	for {
		quote := getQuote(stock.Symbol)

		oldPrice := stock.Price
		stock.Price = quote["c"].(float64)

		if stock.Price > stock.PriceTarget.Float64 && stock.PriceTarget.Valid {
			fmt.Printf("This is the price %v, this is the priceTarget %v", stock.Price, stock.PriceTarget)
			em.PriceTargetMet(fmt.Sprintf(`The following stock has met its
							price target of %v: %v`, stock.PriceTarget, stock.Symbol))
		}

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

func getPriceTarget(reader *bufio.Reader) float64 {
	priceString, _ := reader.ReadString('\n')
	priceString = strings.TrimSuffix(priceString, "\n")
	priceTarget, _ := strconv.ParseFloat(priceString, 64)

	return priceTarget
}

func watch(db *sql.DB) {
	rows, err := db.Query("SELECT symbol, name, price, pricetarget FROM stocks")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var symbol string
		var name sql.NullString
		var price float64
		var priceTarget sql.NullFloat64

		err = rows.Scan(&symbol, &name, &price, &priceTarget)
		if err != nil {
			panic(err)
		}
		s := Stock{
			Symbol:      symbol,
			Name:        name,
			Price:       price,
			PriceTarget: priceTarget,
		}

		go watchStock(&s, db)
	}
}

func insert(db *sql.DB) {
	fmt.Println("What is the symbol of the stock you want to add?")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter symbol: ")
	symbol, _ := reader.ReadString('\n')
	symbol = strings.TrimSuffix(symbol, "\n")
	fmt.Println("What is the name of the stock you want to add?")
	fmt.Print("Enter name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSuffix(name, "\n")
	var valid bool
	if name != "" {
		valid = true
	} else {
		valid = false
	}
	fmt.Println("What is the price target of the stock you want to add?")
	fmt.Print("Enter price target: ")
	priceTarget := getPriceTarget(reader)
	s := Stock{
		Symbol: symbol,
		Name: sql.NullString{
			String: name,
			Valid:  valid,
		},
		PriceTarget: sql.NullFloat64{
			Float64: priceTarget,
			Valid:   true,
		},
	}

	s.insertStock(db)
}

func update(db *sql.DB) {
	fmt.Println("What is the symbol of the stock you want to update?")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter symbol: ")
	symbol, _ := reader.ReadString('\n')
	symbol = strings.TrimSuffix(symbol, "\n")
	fmt.Printf("What is the new price target you want to update %v to? ", symbol)
	priceTarget := getPriceTarget(reader)

	updatePriceTarget(symbol, priceTarget, db)
}

func remove(db *sql.DB) {
	fmt.Println("What is the symbol of the stock you want to remove?")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter symbol: ")
	symbol, _ := reader.ReadString('\n')
	symbol = strings.TrimSuffix(symbol, "\n")

	removeStock(symbol, db)
}

func printUsage() {
	fmt.Print("Usage: ")
	fmt.Println("./main <action>")
	fmt.Println("action == watch will begin a server that watches the stocks that are in the database")
	fmt.Println("action == insert will insert a new stock into the database")
	fmt.Println("action == update will update a stock's price target to a new value")
	fmt.Println("action == remove will remove a stock from the database")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// db, err := database.Connect(database.Host, database.Port, database.User, database.Password, database.Dbname)
	db, err := database.ConnectHeroku(en.GoDotEnvVariable("HEROKU"))
	if err != nil {
		panic(err)
	}

	if os.Args[1] == "watch" {
		watch(db)
		WaitForCtrlC()
	} else if os.Args[1] == "insert" {
		insert(db)
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Would you like to insert another stock?")
			fmt.Println("yes or y to insert another stock, anything else to quit")
			response, _ := reader.ReadString('\n')
			response = strings.TrimSuffix(response, "\n")
			if response == "yes" || response == "y" {
				insert(db)
			} else {
				fmt.Println("Goodbye!")
				return
			}
		}
	} else if os.Args[1] == "update" { // TODO: Make the for loop a function
		update(db)
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Would you like to update another stock?")
			fmt.Println("yes or y to update another stock, anything else to quit")
			response, _ := reader.ReadString('\n')
			response = strings.TrimSuffix(response, "\n")
			if response == "yes" || response == "y" {
				update(db)
			} else {
				fmt.Println("Goodbye!")
				return
			}
		}
	} else if os.Args[1] == "remove" {
		remove(db)
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Would you like to remove another stock?")
			fmt.Println("yes or y to remove another stock, anything else to quit")
			response, _ := reader.ReadString('\n')
			response = strings.TrimSuffix(response, "\n")
			if response == "yes" || response == "y" {
				remove(db)
			} else {
				fmt.Println("Goodbye!")
				return
			}
		}
	} else {
		printUsage()
	}
}
