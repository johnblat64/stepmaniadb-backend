package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	sqlBuilder "github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackskj/carta"
	"go.uber.org/zap"
)

var gDbConnection *pgx.Conn
var gLogger *zap.Logger
var gSugar *zap.SugaredLogger

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

func getListOfSongsFromSearchParameters(c *gin.Context) {

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

type example struct {
	Name     string  `json:"name"`
	Number   int     `json:"number"`
	IsCool   bool    `json:"isCool"`
	PriceMax float32 `json:"priceMax"`
	PriceMin float32 `json:"priceMin"`
}

var examples = []example{
	{Name: "Booger", Number: 1, IsCool: true, PriceMin: 5, PriceMax: 10},
}

func postExample(c *gin.Context) {
	var newExample example

	if err := c.BindJSON(&newExample); err != nil {
		return
	}

	examples = append(examples, newExample)
	c.IndentedJSON(http.StatusCreated, newExample)
}

// var validQueryParams = []string{
// 	"name",
// 	"number",
// 	"isCool",
// 	"priceMax",
// 	"priceMin",
// }

var queryParamTextFilters = []string{
	"songid",
	"packid",
	"packname",
	"title",
	"version",
	"subtitle",
	"artist",
	"genre",
	"songcategory",
}

const DEFAULT_COUNT = "20"

type TimeSignature struct {
	Numerator   int `json:"numerator"  db:"song.timesignature.numerator"`
	Denominator int `json:"denominator" db:"song.timesignature.denominator"`
}

type Chart struct {
	ChartId      string `json:"chartId" db:"chart.chartid"`
	ChartName    string `json:"name" db:"chart.chartname"`
	StepsType    string `json:"stepsType" db:"chart.stepstype"`
	Description  string `json:"description" db:"chart.description"`
	ChartStyle   string `json:"chartStyle" db:"chart.chartstyle"`
	Difficulty   string `json:"difficulty" db:"chart.difficulty"`
	Meter        int    `json:"meter" db:"chart.meter"`
	Credit       string `json:"credit" db:"chart.credit"`
	StopsCount   int    `json:"stopsCount" db:"chart.stops_count"`
	DelaysCount  int    `json:"delaysCount" db:"chart.delays_count"`
	WarpsCount   int    `json:"warpsCount" db:"chart.warps_count"`
	ScrollsCount int    `json:"scrollsCount" db:"chart.scrolls_count"`
	FakesCount   int    `json:"fakesCount" db:"chart.fakes_count"`
	SpeedsCount  int    `json:"speedsCount" db:"chart.speeds_count"`
}

type Song struct {
	SongId         string          `json:"songId"         db:"song.songid"`
	SongTitle      string          `json:"songTitle"      db:"song.title"`
	SongArtist     string          `json:"songArtist"     db:"song.artist"`
	PackId         string          `json:"packId"         db:"pack.packid"`
	PackName       string          `json:"packName"       db:"pack.name"`
	Bpms           []float32       `json:"bpms"           db:"song.bpms"`
	TimeSignatures []TimeSignature `json:"timeSignatures" db:"song.timesignatures"`
	Charts         []Chart
}

// https://github.com/jackskj/carta

var songEntryStruct = sqlBuilder.NewStruct(new(Song))

func getSongList(c *gin.Context) {
	sqlSelectBuilder := sqlBuilder.NewSelectBuilder()

	sqlSelectBuilder.Select(
		"song.songid",
		"song.title",
		"song.artist",
		"pack.packid",
		"pack.name",
		"song.bpms",
		"song.timesignatures",
		"chart.chartid",
		"chart.meter")

	// "song.version",
	// "song.subtitle",
	// "song.genre",
	// "song.songcategory",)

	sqlSelectBuilder.From("song")
	sqlSelectBuilder.Join("pack_song_map", "pack_song_map.songid = song.songid")
	sqlSelectBuilder.Join("pack", "pack.packid = pack_song_map.packid")
	sqlSelectBuilder.Join("chart", "chart.songid = song.songid")
	sqlSelectBuilder.GroupBy("song.songid", "pack.packid")
	countStr := c.DefaultQuery("count", DEFAULT_COUNT)
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return
	}
	sqlSelectBuilder.Limit(count)
	for _, queryParam := range queryParamTextFilters {
		if value, exists := c.GetQuery(queryParam); exists {
			sqlSelectBuilder.Where(strings.ToLower(queryParam) + "='" + value + "'")
		}
	}

	sqlQuery, _ := sqlSelectBuilder.BuildWithFlavor(sqlBuilder.PostgreSQL)
	rows, err := gDbConnection.Query(context.Background(), sqlQuery)
	if err != nil {
		gSugar.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	songs := []Song{}
	carta.Map(rows, &songs)
	var songEntries []Song
	if rows.Err() != nil {
		gSugar.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for rows.Next() {
		var songEntry Song
		err := rows.Scan(songEntryStruct.Addr(&songEntry)...)
		if err != nil {
			gSugar.Error(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		songEntries = append(songEntries, songEntry)
	}

	c.IndentedJSON(http.StatusOK, songEntries)

}

func getExampleList(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, examples)
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

	gDbConnection, err = pgx.Connect(context.Background(),
		"user="+dbUser+" dbname="+dbName+" password="+dbPass+" host="+dbHost)
	if err != nil {
		log.Fatal(err)
	}
	gSugar.Info("Connection to DB " + dbName + " established!")

	var ts TimeSignature
	var _ts []TimeSignature
	pgtype.NewMap().RegisterDefaultPgType(ts, "timesignature")
	pgtype.NewMap().RegisterDefaultPgType(_ts, "_timesignature")

	timeSignaturePgType, err := gDbConnection.LoadType(context.Background(), "timesignature")
	if err != nil {
		gSugar.Panic(err)
	}

	gDbConnection.TypeMap().RegisterType(timeSignaturePgType)

	timeSignaturePgTypeArray, err := gDbConnection.LoadType(context.Background(), "_timesignature")
	if err != nil {
		gSugar.Panic(err)
	}

	gDbConnection.TypeMap().RegisterType(timeSignaturePgTypeArray)

	router := gin.Default()
	// router.GET("/exampleList", getExampleList)
	// router.POST("example", postExample)
	router.GET("songList", getSongList)
	router.Run()
}
