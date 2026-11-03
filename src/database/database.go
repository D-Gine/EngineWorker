package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/D-Gine/EngineWorker/src/debug"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDatabase() error {
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PWD")

	debug.PrintMsg("Trying to connect to database...")
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)
	DB, errSql := sql.Open("postgres", psqlSetup)
	if errSql != nil {
		debug.PrintMsg("There is an error while connecting to the database ", errSql)
		return errSql
	}
	debug.PrintMsg("Setup OK. Trying to ping database...")
	err := DB.Ping()
	if err != nil {
		return err
	}
	debug.PrintSuccess("Successfully connected to database!")
	return nil
}
