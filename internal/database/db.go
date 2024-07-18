package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func InitDB() *sql.DB {
	connUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sql.Open("pgx", connUrl)
	if err != nil {
		log.Fatal("failed to connect to the database: ", err)
	}
	return db
}
