package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/MicahAsowata/Galbium/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type application struct {
	Logger         *zap.Logger
	Queries        *models.Queries
	Users          *models.Users
	SessionManager *scs.SessionManager
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		logger.Error(err.Error())
	}
	err = godotenv.Load()
	if err != nil {
		logger.Error(err.Error())
	}

	dbUserName := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dsn := dbUserName + ":" + dbPassword + "@" + dbHost + "/" + dbName + "?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Error(err.Error())
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = time.Hour * 12
	sessionManager.Cookie.Secure = true
	sessionManager.Cookie.Persist = false
	a := application{
		Logger:         logger,
		Queries:        &models.Queries{DB: db},
		Users:          &models.Users{DB: db},
		SessionManager: sessionManager,
	}
	port := ":" + os.Getenv("PORT")
	srv := &http.Server{
		Addr:    port,
		Handler: a.routes(),
	}
	logger.Info("Starting server for http://localhost" + port)
	err = srv.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
	}
}
