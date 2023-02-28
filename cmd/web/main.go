package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/MicahAsowata/Galbium/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	Queries *models.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("could not load the .env file")
	}
	dsn := "galbius:galbius@/galbius?parseTime=true"
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
