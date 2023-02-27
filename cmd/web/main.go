package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/MicahAsowata/Galbium/internal/models"
	"github.com/joho/godotenv"
)

type application struct {
	Queries *models.Queries
}

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("could not load the .env file")
	}
	dsn := os.Getenv("DB_NAME")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	queries := models.New(db)
	a := application{
		Queries: queries,
	}
	log.Println("Starting server for http://localhost:4000")
	err = http.ListenAndServe(":4000", a.routes())
	if err != nil {
		log.Fatal(err)
	}
}
