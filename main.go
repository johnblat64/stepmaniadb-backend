package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib" // needed for sql.Open
	"go.uber.org/zap"
)

var gDb *sql.DB
var gLogger *zap.Logger
var gSugar *zap.SugaredLogger

type DbUserConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
}

type DbEngineType string

const (
	DbEgninePg      = "pg"
	DbEngineSqlite3 = "sqlite3"
)

func main() {
	var err error

	if gLogger, err = zap.NewProduction(); err != nil {
		log.Fatal(err)
	}
	defer gLogger.Sync()

	gSugar = gLogger.Sugar()
	gSugar.Info("Zap Sugared Logger Initialized")

	var dbEngine DbEngineType
	var dbUserConfig DbUserConfig
	var exists bool

	awsRegion, awsRegionEnvVarExists := os.LookupEnv("AWS_REGION")

	awsDbPasswordSecretName, awsDbPasswordSecretNameEnvVarExists := os.LookupEnv("AWS_DB_PASSWORD_SECRET_NAME")

	if awsRegionEnvVarExists && awsDbPasswordSecretNameEnvVarExists {
		awsDbUserSecret := getAwsSecret(awsRegion, awsDbPasswordSecretName)
		dbUserConfig.Username = awsDbUserSecret.Username
		dbUserConfig.Password = awsDbUserSecret.Password
		dbUserConfig.Host = awsDbUserSecret.Host
		dbUserConfig.Port = strconv.FormatInt(int64(awsDbUserSecret.Port), 10)
	} else {
		awsEnvVarMessage := " or AWS_REGION and AWS_DB_PASSWORD_SECRET_NAME in environment"
		dbUserConfig.Password, exists = os.LookupEnv("DB_PASS")
		if !exists {
			log.Fatal("Must have DB_PASS in environment" + awsEnvVarMessage)
		}

		dbUserConfig.Username, exists = os.LookupEnv("DB_USER")
		if !exists {
			log.Fatal("Must have DB_USER in environment" + awsEnvVarMessage)
		}

		dbUserConfig.Host, exists = os.LookupEnv("DB_HOST")
		if !exists {
			log.Fatal("Must have DB_HOST in environment" + awsEnvVarMessage)
		}

		dbUserConfig.Port, exists = os.LookupEnv("DB_PORT")
		if !exists {
			dbUserConfig.Port = "5432"
		}
	}

	dbUserConfig.Name, exists = os.LookupEnv("DB_NAME")
	if !exists {
		log.Fatal("Must have DB_NAME in environment")
	}

	gSugar.Info("Connecting to DB " + dbUserConfig.Name + "...")

	connStr := "user=" + dbUserConfig.Username + " dbname=" + dbUserConfig.Name + " password=" + dbUserConfig.Password + " host=" + dbUserConfig.Host + " port=" + dbUserConfig.Port
	gDb, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer gDb.Close()

	gSugar.Info("Connection to DB " + dbUserConfig.Name + " established!")

	router := gin.Default()

	router.GET("songs", getFilteredSongs)
	router.GET("songs/:songid", getSongById)
	router.GET("packs/:packid", getPackById)
	router.Run()
}
