package client

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/qwuiemme/ellipsespace-server/config"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", config.New().DBConnectionString)

	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(10 * time.Second)
	db.SetMaxOpenConns(100)

	return db
}
