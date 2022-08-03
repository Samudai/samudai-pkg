package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Samudai/samudai-pkg/logger"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitSQL() {
	dbURL := os.Getenv("DATABASE_URL")

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	db = conn
	logger.LogMessage("info", "postgres db connected")

	if err := db.Ping(); err != nil {
		logger.LogMessage("error", "error pinging to db: %v", err)
		logger.LogMessage("debug", "reconnecting")
		InitSQL()
	}
}

func GetSQL() *sql.DB {
	return db
}

func CloseSQL() {
	db.Close()
}
