package client

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "admin:admin@/ellipsespace")

	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(10 * time.Second)
	db.SetMaxOpenConns(100)

	return db
}
