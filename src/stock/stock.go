package stock

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func (s *Stock) UpdateStock(db *sql.DB, newValue float64) {
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

func (s *Stock) InsertStock(db *sql.DB) {
	quote := GetQuote(s.Symbol)

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

func UpdatePriceTarget(symbol string, newPriceTaret float64, db *sql.DB) {
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

func RemoveStock(symbol string, db *sql.DB) {
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

func WatchStock(stock *Stock, db *sql.DB) {
	fmt.Println("Beginning to watch stock", stock.Symbol)
	for {
		quote := GetQuote(stock.Symbol)

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
				stock.UpdateStock(db, newPrice)
			}
		default:
			panic(err)
		}

		time.Sleep(1 * time.Minute)
	}
}

func GetPriceTarget(reader *bufio.Reader) float64 {
	priceString, _ := reader.ReadString('\n')
	priceString = strings.TrimSuffix(priceString, "\n")
	priceTarget, _ := strconv.ParseFloat(priceString, 64)

	return priceTarget
}

func GetQuote(symbol string) map[string]interface{} {
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
