package main

import (
	"context"
	"database/sql"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

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

const DEFAULT_COUNT = "20"

func addWhereClausesFromQueryParameters(sb *sqlBuilder.SelectBuilder, c *gin.Context) {
	// game := c.Query("stepstype")
	// if game != "" {
	// 	sb.Where(sb.Like("chart.stepstype", game+"%"))
	// }

	stepsType := c.Query("stepstype")
	if stepsType != "" {
		sb.Where(sb.E("lower(chart.stepstype)", strings.ToLower(stepsType)))
	}

	meterMin := c.DefaultQuery("meterMin", "0")
	meterMax := c.DefaultQuery("meterMax", "99")
	sb.Where(sb.Between("chart.meter", meterMin, meterMax))

	bpmMin := c.DefaultQuery("bpmMin", "0")
	bpmMax := c.DefaultQuery("bpmMax", "999")

	sb.Where(sb.Between("song_bpm.song_bpm", bpmMin, bpmMax))

	packTitle := c.Query("pack")
	if packTitle != "" {
		sb.Where(sb.Like("lower(pack.title)", "%"+strings.ToLower(packTitle)+"%"))
	}

	timeSignatureNumerator := c.Query("timeSignatureNumerator")
	if timeSignatureNumerator != "" {
		sb.Where(sb.E("song_time_signature.time_signature_numerator", timeSignatureNumerator))
	}

	timeSignatureDenominator := c.Query("timeSignatureDenominator")
	if timeSignatureDenominator != "" {
		sb.Where(sb.E("song_time_signature.time_signature_denominator", timeSignatureDenominator))
	}

	songTitle := c.Query("title")
	if songTitle != "" {
		sb.Where(sb.Like("lower(song.title)", "%"+strings.ToLower(songTitle)+"%"))
	}

	songArtist := c.Query("artist")
	if songArtist != "" {
		sb.Where(sb.Like("lower(song.artist)", "%"+strings.ToLower(songArtist)+"%"))
	}

	chartDifficultyMeterMin := c.Query("chartDifficultyMeterMin")
	chartDifficultyMeterMax := c.Query("chartDifficultyMeterMax")
	if chartDifficultyMeterMin != "" && chartDifficultyMeterMax != "" {
		sb.Where(sb.Between("chart.meter", chartDifficultyMeterMin, chartDifficultyMeterMax))
	}

	chartCredit := c.Query("chartCredit")
	if chartCredit != "" {
		sb.Where(sb.Like("lower(chart.credit)", "%"+strings.ToLower(chartCredit)+"%"))
	}
}

func getFilteredSongs(c *gin.Context) {
	sbFilteredSongs := sqlBuilder.NewSelectBuilder()

	// get pageStr from query param
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		gSugar.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	DEFAULT_COUNT := "20"
	MAX_COUNT := 100
	// get pageSizeStr from query param
	pageSizeStr := c.DefaultQuery("pageSize", DEFAULT_COUNT)
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		gSugar.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if pageSize > MAX_COUNT {
		pageSize = MAX_COUNT
	}

	// Count of songs in subquery
	sbCountFilteredSongs := sqlBuilder.NewSelectBuilder()
	sbCountFilteredSongs.Select("COUNT( DISTINCT song.songid)").As("count", "count")
	sbCountFilteredSongs.From("song")
	sbCountFilteredSongs.Join("pack_song_map", "pack_song_map.songid = song.songid")
	sbCountFilteredSongs.Join("pack", "pack.packid = pack_song_map.packid")
	sbCountFilteredSongs.Join("chart", "chart.songid = song.songid")
	sbCountFilteredSongs.Join("song_bpm", "song_bpm.songid = song.songid")
	sbCountFilteredSongs.Join("song_time_signature", "song_time_signature.songid = song.songid")
	addWhereClausesFromQueryParameters(sbCountFilteredSongs, c)
	countQuery, args := sbCountFilteredSongs.BuildWithFlavor(sqlBuilder.PostgreSQL)
	// get filteredSongsCount value in filteredSongsCount struct
	var filteredSongsCount Count
	countRow := gDb.QueryRow(countQuery, args...)
	err = countRow.Scan(&filteredSongsCount.Count)
	if err != nil {
		gSugar.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	offset := (page - 1) * pageSize
	limit := pageSize
	pageCount := int(math.Ceil(float64(filteredSongsCount.Count) / float64(pageSize)))
	totalSongsCount := filteredSongsCount.Count

	// Subquery
	sbFilteredSongs.Select("song.songid, song.title").Distinct().As("song.songid", "songid")
	sbFilteredSongs.From("song")
	sbFilteredSongs.Join("pack_song_map", "pack_song_map.songid = song.songid")
	sbFilteredSongs.Join("pack", "pack.packid = pack_song_map.packid")
	sbFilteredSongs.Join("chart", "chart.songid = song.songid")
	sbFilteredSongs.Join("song_bpm", "song_bpm.songid = song.songid")
	sbFilteredSongs.Join("song_time_signature", "song_time_signature.songid = song.songid")
	sbFilteredSongs.OrderBy("song.title ASC")
	sbFilteredSongs.Offset(offset)
	sbFilteredSongs.Limit(limit)
	addWhereClausesFromQueryParameters(sbFilteredSongs, c)

	// Main Query
	sbSongData := sqlBuilder.NewSelectBuilder()
	sbSongData.Select("song.songid, song.title, song.artist, song.banner_path, song.music_path, song.file_extension, song.song_dir_path, song_bpm.song_bpm, song_time_signature.time_signature_numerator, song_time_signature.time_signature_denominator, chart.chartid, chart.chartname, chart.stepstype, chart.description, chart.chartstyle, chart.difficulty, chart.meter, chart.credit, chart.stops_count, chart.delays_count, chart.warps_count, chart.scrolls_count, chart.fakes_count, chart.speeds_count, chart.stream, chart.voltage, chart.air, chart.freeze, chart.chaos, pack.packid, pack.name")
	sbSongData.From(sbSongData.BuilderAs(sbFilteredSongs, "filtered_songs"))
	sbSongData.Join("song", "filtered_songs.songid = song.songid")
	sbSongData.Join("song_bpm", "song_bpm.songid = song.songid")
	sbSongData.Join("song_time_signature", "song_time_signature.songid = song.songid")
	sbSongData.Join("chart", "chart.songid = song.songid")
	sbSongData.Join("pack_song_map", "pack_song_map.songid = song.songid")
	sbSongData.Join("pack", "pack.packid = pack_song_map.packid")
	sbSongData.OrderBy("song.title ASC")

	songDataQuery, args := sbSongData.BuildWithFlavor(sqlBuilder.PostgreSQL)
	gSugar.Info(songDataQuery)
	gSugar.Info("Args: ", args)

	rows, err := gDb.Query(songDataQuery, args...)
	if err != nil {
		gSugar.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	sbCount := sqlBuilder.NewSelectBuilder()
	sbCount.Select("COUNT(songid)")
	sbCount.From("song")
	songResultsResponse := SongResultsResponse{}
	songResultsResponse.PageCount = pageCount
	songResultsResponse.Page = page
	songResultsResponse.PageSize = pageSize
	songResultsResponse.TotalSongsCount = totalSongsCount
	err = carta.Map(rows, &songResultsResponse.Songs)
	// songPage.ResultsCount = len(songPage.Songs)
	if err != nil {
		gSugar.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, songResultsResponse)
}

func getSongById(c *gin.Context) {
	songid := c.Param("songid")

	sqlSongSelectBuilder := sqlBuilder.NewSelectBuilder()
	sqlSongSelectBuilder.Select("song.songid, song.title, song.artist, song.banner_path, song.music_path, song.song_dir_path, song_bpm.song_bpm, song_time_signature.time_signature_numerator, song_time_signature.time_signature_denominator, chart.chartid, chart.chartname, chart.stepstype,  chart.description, chart.chartstyle,   chart.difficulty, chart.meter, chart.credit, chart.stream, chart.voltage, chart.air, chart.freeze, chart.chaos, pack.packid, pack.name ")
	sqlSongSelectBuilder.From("song")
	sqlSongSelectBuilder.Join("pack_song_map", "pack_song_map.songid = song.songid")
	sqlSongSelectBuilder.Join("pack", "pack.packid = pack_song_map.packid")
	sqlSongSelectBuilder.Join("chart", "chart.songid = song.songid")
	sqlSongSelectBuilder.Join("song_bpm", "song_bpm.songid = song.songid")
	sqlSongSelectBuilder.Join("song_time_signature", "song_time_signature.songid = song.songid")
	sqlSongSelectBuilder.Where(sqlSongSelectBuilder.E("song.songid", songid))
	sqlSongQuery, args := sqlSongSelectBuilder.BuildWithFlavor(sqlBuilder.PostgreSQL)

	rows, err := gDb.Query(sqlSongQuery, args...)
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

	c.JSON(200, songs[0])
}

func getPackById(c *gin.Context) {
	packid := c.Param("packid")

	sqlPackSelectBuilder := sqlBuilder.NewSelectBuilder()
	sqlPackSelectBuilder.Select("song.songid, song.title, song.artist,  song.banner_path, song.music_path, song.song_dir_path, song_bpm.song_bpm, song_time_signature.time_signature_numerator, song_time_signature.time_signature_denominator, chart.meter, pack.packid, pack.name, pack.download_link ")
	sqlPackSelectBuilder.From("song")
	sqlPackSelectBuilder.Join("pack_song_map", "pack_song_map.songid = song.songid")
	sqlPackSelectBuilder.Join("pack", "pack.packid = pack_song_map.packid")
	sqlPackSelectBuilder.Join("chart", "chart.songid = song.songid")
	sqlPackSelectBuilder.Join("song_bpm", "song_bpm.songid = song.songid")
	sqlPackSelectBuilder.Join("song_time_signature", "song_time_signature.songid = song.songid")
	sqlPackSelectBuilder.OrderBy("song.title ASC")
	sqlPackSelectBuilder.Where(sqlPackSelectBuilder.E("pack.packid", packid))
	sqlPackQuery, args := sqlPackSelectBuilder.BuildWithFlavor(sqlBuilder.PostgreSQL)

	rows, err := gDb.Query(sqlPackQuery, args...)
	defer rows.Close()
	if err != nil {
		gSugar.Errorln(err)
		c.Data(500, "application/text", []byte("Internal Server Error"))
	}

	packs := []Pack{}

	err = carta.Map(rows, &packs)
	if err != nil {
		gSugar.Errorln(err)
	}

	c.JSON(200, packs[0])
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
	router.GET("packs/:packid", getPackById)
	router.Run()
}
