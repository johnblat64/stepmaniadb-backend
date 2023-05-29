package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib" // needed for sql.Open
	"go.uber.org/zap"
)

var gDb *sql.DB
var gLogger *zap.Logger
var gSugar *zap.SugaredLogger

func main() {
	var err error

	if gLogger, err = zap.NewProduction(); err != nil {
		log.Fatal(err)
	}
	defer gLogger.Sync()

	gSugar = gLogger.Sugar()
	gSugar.Info("Zap Sugared Logger Initialized")

	dbPass, exists := os.LookupEnv("DB_PASS")
	if !exists {
		log.Fatal("Must have DB_PASS in environment")
	}
	dbUser, exists := os.LookupEnv("DB_USER")
	if !exists {
		log.Fatal("Must have DB_USER in environment")
	}
	dbName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		log.Fatal("Must have DB_NAME in environment")
	}
	dbHost, exists := os.LookupEnv("DB_HOST")
	if !exists {
		log.Fatal("Must have DB_HOST in environment")
	}

	gSugar.Info("Connecting to DB " + dbName + "...")

	connStr := "user=" + dbUser + " dbname=" + dbName + " password=" + dbPass + " host=" + dbHost
	gDb, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer gDb.Close()

	gSugar.Info("Connection to DB " + dbName + " established!")

	router := gin.Default()

	router.GET("songs", getFilteredSongs)
	router.GET("songs/:songid", getSongById)
	router.GET("packs/:packid", getPackById)
	router.Run()
}
