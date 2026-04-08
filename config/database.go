package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {

	connStr := fmt.Sprintf(
		"host=localhost user=postgres password=2863 dbname=weavory_db port=5432 sslmode=disable",
	)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Database Connected")

	DB = db
}