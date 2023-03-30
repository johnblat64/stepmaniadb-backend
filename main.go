package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	sqlBuilder "github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

type Whatever struct {
	X []byte `db:"song_chart_json"`
}

const DEFAULT_COUNT = "20"

func getFilteredSongs(c *gin.Context) {
	sqlSelectBuilder := sqlBuilder.NewSelectBuilder()

	sqlSelectBuilder.Select("song_chart_json")
	sqlSelectBuilder.From("filtered_songs_with_charts('',0,999,4,4,'','pump-single',	0, 99, '')")

	sqlQuery, _ := sqlSelectBuilder.BuildWithFlavor(sqlBuilder.PostgreSQL)

	rows, err := gDbConnection.Query(context.Background(), sqlQuery)
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
	// json.Marshal(whatever.X)
	c.Data(http.StatusOK, "application/json", whatever.X)
	//c.IndentedJSON()
}

type SongsResponse struct {
	SongId         string          `json:"songId"         db:"song.songid"`
	SongTitle      string          `json:"songTitle"      db:"song.title"`
	SongArtist     string          `json:"songArtist"     db:"song.artist"`
	PackId         string          `json:"packId"         db:"pack.packid"`
	PackName       string          `json:"packName"       db:"pack.name"`
	Bpms           []float32       `json:"bpms"           db:"song.bpms"`
	TimeSignatures []TimeSignature `json:"timeSignatures" db:"song.timesignatures"`
	Charts         []Chart         `json:"charts"`
}

func getSongById(c *gin.Context) {
	songid := c.Param("songid")

	sqlSongSelectBuilder := sqlBuilder.NewSelectBuilder()
	sqlSongSelectBuilder.Select("song.songid, song.title, song.artist, pack.packid, pack.name, song.bpms, song.timesignatures")
	sqlSongSelectBuilder.From("song")
	sqlSongSelectBuilder.Join("pack_song_map", "pack_song_map.songid = song.songid")
	sqlSongSelectBuilder.Join("pack", "pack.packid = pack_song_map.packid")
	sqlSongSelectBuilder.Where("song.songid='" + songid + "'")
	sqlSongQuery, _ := sqlSongSelectBuilder.BuildWithFlavor(sqlBuilder.PostgreSQL)

	var songRow SongRow
	var charts []Chart
	songStruct := sqlBuilder.NewStruct(new(SongRow))

	row := gDbConnection.QueryRow(context.Background(), sqlSongQuery)

	err := row.Scan(songStruct.Addr(&songRow)...)
	if err != nil {
		log.Println(err)
	}

	sqlChartSelectBuilder := sqlBuilder.NewSelectBuilder()
	sqlChartSelectBuilder.Select("chartid, chartname, stepstype,  description, chartstyle,   difficulty, meter, credit, stops_count, delays_count, warps_count, scrolls_count, fakes_count, speeds_count")
	sqlChartSelectBuilder.From("chart")
	sqlChartSelectBuilder.Where("chart.songid='" + songid + "'")
	sqlChartQuery, _ := sqlChartSelectBuilder.BuildWithFlavor(sqlBuilder.PostgreSQL)

	rows, err := gDbConnection.Query(context.Background(), sqlChartQuery)
	defer rows.Close()
	if rows.Err() != nil {
		log.Println(rows.Err())
	}
	for rows.Next() {
		var chart Chart
		chartStruct := sqlBuilder.NewStruct(new(Chart))
		err := rows.Scan(chartStruct.Addr(&chart)...)
		if rows.Err() != nil {
			log.Println(rows.Err())
		}
		if err != nil {
			log.Println(err)
		}
		charts = append(charts, chart)
	}

	songResponse := SongsResponse{
		SongId:         songRow.SongId,
		SongTitle:      songRow.SongTitle,
		SongArtist:     songRow.SongArtist,
		PackId:         songRow.PackId,
		PackName:       songRow.PackName,
		Bpms:           songRow.Bpms,
		TimeSignatures: songRow.TimeSignatures,
		Charts:         charts,
	}
	c.JSON(200, songResponse)
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

	gDbConnection, err = pgx.Connect(context.Background(), "user="+dbUser+" dbname="+dbName+" password="+dbPass+" host="+dbHost)
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
	router.GET("songList", getFilteredSongs)
	router.GET("song/:songid", getSongById)
	router.Run()
}
