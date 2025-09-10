package db

import (
	"database/sql"
	"log"
	"os"
	"fmt"

	"github.com/joho/godotenv"
	 _ "github.com/lib/pq"
)

var DB *sql.DB

func Connect()  {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

	DB, err = sql.Open("postgres", psqlInfo) 
	if err != nil {
		log.Fatal("Error opening DB:", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal("Error pinging db: ",err)
	}
	fmt.Println("yeeee......Connected to PostgreSQL")
}