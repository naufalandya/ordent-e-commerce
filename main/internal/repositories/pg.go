package repositories

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

var DB *pgx.Conn

func InitDB() {
	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	fmt.Println("DB_USER:", dbUser)
	fmt.Println("DB_PASSWORD:", dbPassword)
	fmt.Println("DB_HOST:", dbHost)
	fmt.Println("DB_PORT:", dbPort)
	fmt.Println("DB_NAME:", dbName)

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", dbUser, dbPassword, dbHost, dbPort, dbName)

	DB, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	fmt.Println("Connected to PostgreSQL")
}
