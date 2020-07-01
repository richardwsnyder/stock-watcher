package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	e "github.com/richardwsnyder/stock-watcher/src/env"
)

var (
	Host     = e.GoDotEnvVariable("HOST")
	Port     = 5432
	User     = e.GoDotEnvVariable("USER")
	Password = e.GoDotEnvVariable("PASSWORD")
	Dbname   = e.GoDotEnvVariable("DBNAME")
)

func Connect(host string, port int, user string, password string, dbname string) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}
