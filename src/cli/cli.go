package cli

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	st "github.com/richardwsnyder/stock-watcher/src/stock"
)

func Watch(db *sql.DB) {
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
		s := st.Stock{
			Symbol:      symbol,
			Name:        name,
			Price:       price,
			PriceTarget: priceTarget,
		}

		go st.WatchStock(&s, db)
	}
}

func Insert(db *sql.DB) {
	fmt.Println("What is the symbol of the stock you want to add?")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter symbol: ")
	symbol, _ := reader.ReadString('\n')
	symbol = strings.ToUpper(strings.TrimSuffix(symbol, "\n"))
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
	priceTarget := st.GetPriceTarget(reader)
	s := st.Stock{
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

	s.InsertStock(db)
}

func Update(db *sql.DB) {
	fmt.Println("What is the symbol of the stock you want to update?")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter symbol: ")
	symbol, _ := reader.ReadString('\n')
	symbol = strings.ToUpper(strings.TrimSuffix(symbol, "\n"))
	fmt.Printf("What is the new price target you want to update %v to? ", symbol)
	priceTarget := st.GetPriceTarget(reader)

	st.UpdatePriceTarget(symbol, priceTarget, db)
}

func Remove(db *sql.DB) {
	fmt.Println("What is the symbol of the stock you want to remove?")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter symbol: ")
	symbol, _ := reader.ReadString('\n')
	symbol = strings.ToUpper(strings.TrimSuffix(symbol, "\n"))

	st.RemoveStock(symbol, db)
}
