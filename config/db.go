package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DbConn struct {
	Db *sql.DB
}

func ConnectDB() *DbConn {
	DB, err := sql.Open("mysql", "myuser:mypassword@tcp(localhost:3307)/mydatabase")
	if err != nil {
		log.Fatal("cannot estabilish connection to mysql dtabase", err)
	}
	return &DbConn{Db: DB}

}
