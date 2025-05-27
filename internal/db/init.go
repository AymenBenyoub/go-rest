package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// InitDB initializes the database connection and returns a *sql.DB instance.

func InitDB() (*sql.DB, error) {

	db, err := sql.Open("mysql", "go_user:kappa123@tcp(localhost:3306)/go_test_db")
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {

		return nil, err
	}
	
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)
	return db, nil
}
