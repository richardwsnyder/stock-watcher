package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

var (
	Host     = GoDotEnvVariable("HOST")
	Port     = 5432
	User     = GoDotEnvVariable("USER")
	Password = GoDotEnvVariable("PASSWORD")
	Dbname   = GoDotEnvVariable("DBNAME")
	Token    = GoDotEnvVariable("TOKEN")
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

// func main() {
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
// 		"password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbname)
// 	db, err := sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	sqlStatement := `
// 	INSERT INTO users (age, email, first_name, last_name)
// 	VALUES ($1, $2, $3, $4)
// 	RETURNING id`
// 	id := 0
// 	err = db.QueryRow(sqlStatement, 30, "jon@calhoun.io", "Jonathan", "Calhoun").Scan(&id)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("New record ID is:", id)
// }
