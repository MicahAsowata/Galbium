package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
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
	// dsn := "galbius:galbius@/galbius?parseTime=true"
	dbUserName := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dsn := dbUserName + ":" + dbPassword + "@" + dbHost + "/" + dbName + "?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = time.Hour * 12
	sessionManager.Cookie.Persist = false
	a := application{
		Queries:        &models.Queries{DB: db},
		Users:          &models.Users{DB: db},
		SessionManager: sessionManager,
	}
	port := ":" + os.Getenv("PORT")
	log.Println("Starting server for http://localhost" + port)
	err = http.ListenAndServe(port, a.routes())
	if err != nil {
		log.Fatal(err)
	}
}
