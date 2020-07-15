package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"

	cl "github.com/richardwsnyder/stock-watcher/src/cli"
	"github.com/richardwsnyder/stock-watcher/src/database"
	en "github.com/richardwsnyder/stock-watcher/src/env"
	se "github.com/richardwsnyder/stock-watcher/src/server"
)

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

func printUsage() {
	fmt.Println("Usage for the command line interface: ")
	fmt.Println(" w will begin a server that watches the stocks that are in the database")
	fmt.Println(" i will insert a new stock into the database")
	fmt.Println(" u will update a stock's price target to a new value")
	fmt.Println(" r will remove a stock from the database")
	fmt.Println(" s will start an http server")
	fmt.Println(" l will list all of the stocks you have currently")
	fmt.Println(" q will quit the program")
}

func getInput(reader *bufio.Reader) string {
	fmt.Println("\nWhat would you like to do?")
	fmt.Println("Type h for help\n")

	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSuffix(response, "\n"))

	return response
}

func main() {
	printUsage()

	// db, err := database.Connect(database.Host, database.Port, database.User, database.Password, database.Dbname)
	db, err := database.ConnectHeroku(en.GoDotEnvVariable("HEROKU"))
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		response := getInput(reader)

		switch response {
		case "h":
			printUsage()
		case "i":
			cl.Insert(db)
		case "l":
			cl.List(db)
		case "q":
			return
		case "r":
			cl.Remove(db)
		case "s":
			go se.Serve(db)
		case "u":
			cl.Update(db)
		case "w":
			cl.Watch(db)
			WaitForCtrlC()
		default:
			printUsage()
		}
	}
}
