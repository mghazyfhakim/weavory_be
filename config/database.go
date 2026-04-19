package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectDB() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found (production mode)")
	}

	dbURL := os.Getenv("DATABASE_URL")

	var db *sql.DB

	if dbURL != "" {
		log.Println("Using DATABASE_URL (production)")

		var err error
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatal("Failed to open DB:", err)
		}

	} else {
		log.Println("Using DB config (local)")

		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")

		connStr := "host=" + host +
			" user=" + user +
			" password=" + password +
			" dbname=" + dbname +
			" port=" + port +
			" sslmode=disable"

		var err error
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal("Failed to open DB:", err)
		}
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect DB:", err)
	}

	log.Println("✅ Database Connected")

	DB = db
}