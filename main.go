package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	sqlBuilder "github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackskj/carta"
	"go.uber.org/zap"
)

var gDb *sql.DB
var gLogger *zap.Logger
var gSugar *zap.SugaredLogger

// ~/go/pkg/mod/github.com/jackc/pgx@v3.6.2+incompatible
// ~/go/pkg/mod/github.com/jackc/pgx/v5@v5.3.1/stdlib
type songSearchParameters struct {
	Game                     string `json:"game"`
	StepsType                string `json:"stepsType"`
	PackTitle                string `json:"packTitle"`
	SongTitle                string `json:"songTitle"`
	SongArtist               string `json:"songArtist"`
	ChartCredit              string `json:"chartCredit"`
	ChartDifficultyMeterMin  int    `json:"chartDifficultyMeterMin"`
	ChartDifficultyMeterMax  int    `json:"chartDifficultyMeterMax"`
	BpmMin                   int    `json:"bpmMin"`
	BpmMax                   int    `json:"bpmMax"`
	TimeSignatureNumerator   int    `json:"timeSignatureNumerator"`
	TimeSignatureDenominator int    `json:"timeSignatureDenominator"`
}

type packAddFromDownloadLinkRequest struct {
	DownloadLink string `json:"downloadLink"`
}

type Whatever struct {
	X []byte `db:"song_chart_json"`
}

const DEFAULT_COUNT = "20"

func getFilteredSongs(c *gin.Context) {
	sqlSelectBuilder := sqlBuilder.NewSelectBuilder()

	sqlSelectBuilder.Select("song_chart_json")
	sqlSelectBuilder.From("filtered_songs_with_charts('',0,999,4,4,'','pump-single',	0, 99, '')")

	sqlQuery, _ := sqlSelectBuilder.BuildWithFlavor(sqlBuilder.PostgreSQL)

	rows, err := gDb.Query(sqlQuery)
	if err != nil {
		gSugar.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if rows.Err() != nil {
		gSugar.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	var whatever Whatever
	var whateverStruct = sqlBuilder.NewStruct(new(Whatever))

	rows.Next()
	if rows.Err() != nil {
		gSugar.Error(rows.Err())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = rows.Scan(whateverStruct.Addr(&whatever)...)
	if err != nil {
		gSugar.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "application/json", whatever.X)
}

func getSongById(c *gin.Context) {
	songid := c.Param("songid")

	sqlSongSelectBuilder := sqlBuilder.NewSelectBuilder()
	sqlSongSelectBuilder.Select("song.songid, song.title, song.artist, song_bpm.song_bpm, song_time_signature.time_signature_numerator, song_time_signature.time_signature_denominator, chart.chartid, chart.chartname, chart.stepstype,  chart.description, chart.chartstyle,   chart.difficulty, chart.meter, chart.credit ")
	sqlSongSelectBuilder.From("song")
	sqlSongSelectBuilder.Join("pack_song_map", "pack_song_map.songid = song.songid")
	sqlSongSelectBuilder.Join("pack", "pack.packid = pack_song_map.packid")
	sqlSongSelectBuilder.Join("chart", "chart.songid = song.songid")
	sqlSongSelectBuilder.Join("song_bpm", "song_bpm.songid = song.songid")
	sqlSongSelectBuilder.Join("song_time_signature", "song_time_signature.songid = song.songid")
	sqlSongSelectBuilder.Where("song.songid='" + songid + "'")
	sqlSongQuery, _ := sqlSongSelectBuilder.BuildWithFlavor(sqlBuilder.PostgreSQL)

	// var songRow SongRow
	// var charts []Chart

	rows, err := gDb.Query(sqlSongQuery)
	defer rows.Close()
	if err != nil {
		gSugar.Errorln(err)
		c.Data(500, "application/text", []byte("Internal Server Error"))
	}

	songs := []Song{}

	err = carta.Map(rows, &songs)
	if err != nil {
		gSugar.Errorln(err)
	}

	c.JSON(200, songs)
}

/**
* Submits a downlaod link to a pack that the user wants to upload to the site.
* As of right now, the download link will be approved by the site admin.
* Once approved, the pack and song metadata will be manually added to the database by the site admin.
* The download link will be available for
* users to download
 */
func postPackAddFromDownloadLinkRequest(c *gin.Context) {

}

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

	// connConfig, err := pgx.ParseConfig("user=" + dbUser + " dbname=" + dbName + " password=" + dbPass + " host=" + dbHost)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	connStr := "user=" + dbUser + " dbname=" + dbName + " password=" + dbPass + " host=" + dbHost
	gDb, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := gDb.Conn(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Raw(func(driverConn any) error {
		pgxConn := driverConn.(*stdlib.Conn).Conn()
		var ts TimeSignature
		var _ts []TimeSignature
		pgtype.NewMap().RegisterDefaultPgType(ts, "timesignature")
		pgtype.NewMap().RegisterDefaultPgType(_ts, "_timesignature")

		timeSignaturePgType, err := pgxConn.LoadType(context.Background(), "timesignature")
		if err != nil {
			gSugar.Panic(err)
		}

		pgxConn.TypeMap().RegisterType(timeSignaturePgType)

		timeSignaturePgTypeArray, err := pgxConn.LoadType(context.Background(), "_timesignature")
		if err != nil {
			gSugar.Panic(err)
		}

		pgxConn.TypeMap().RegisterType(timeSignaturePgTypeArray)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	gSugar.Info("Connection to DB " + dbName + " established!")

	router := gin.Default()
	// router.GET("/exampleList", getExampleList)
	// router.POST("example", postExample)
	router.GET("songs", getFilteredSongs)
	router.GET("songs/:songid", getSongById)
	router.Run()
}
