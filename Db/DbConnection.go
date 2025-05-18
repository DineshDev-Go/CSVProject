package db

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func LocalDBConnect() (*gorm.DB, error) {
	log.Println("LocalDBConnect+")

	// Load .env file (optional, for local dev)
	err := godotenv.Load("config.env")
	if err != nil {
		log.Println("Warning: .env file not loaded. Make sure env variables are set.")
	}

	// Fetch values from environment
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	// Build DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)

	// Connect with GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("GORM connection failed:", err.Error())
		return nil, err
	}

	log.Println("LocalDBConnect-")
	return db, nil
}
