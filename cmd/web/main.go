package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/MicahAsowata/Galbium/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	Queries        *models.Queries
	Users          *models.Users
	SessionManager *scs.SessionManager
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
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = time.Hour * 12
	a := application{
		Queries:        queries,
		Users:          &models.Users{DB: db},
		SessionManager: sessionManager,
	}
	log.Println("Starting server for http://localhost:4000")
	err = http.ListenAndServe(":4000", a.routes())
	if err != nil {
		log.Fatal(err)
	}
}
