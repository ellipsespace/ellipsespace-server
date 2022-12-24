package client

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/qwuiemme/ellipsespace-server/config"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", config.New().DBConnectionString)

	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxIdleTime(0)
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return db
}
