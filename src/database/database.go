package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDatabase() error {
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PWD")

	fmt.Println("Trying to connect to database...")
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)
	DB, errSql := sql.Open("postgres", psqlSetup)
	if errSql != nil {
		fmt.Println("There is an error while connecting to the database ", errSql)
		return errSql
	}
	fmt.Println("Setup OK. Trying to ping database...")
	err := DB.Ping()
	if err != nil {
		return err
	}
	fmt.Println("Ping OK\nSuccessfully connected to database!")
	return nil
}
