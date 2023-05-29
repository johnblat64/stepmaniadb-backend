package main

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	sqlBuilder "github.com/huandu/go-sqlbuilder"
	"github.com/jackskj/carta"
	smdbcore "github.com/stepmaniadb/stepmaniadb-core"
)

const DEFAULT_COUNT = "20"

func addWhereClausesFromQueryParameters(sb *sqlBuilder.SelectBuilder, c *gin.Context) {
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

	packName := c.Query("pack")
	if packName != "" {
		sb.Where(sb.Like("lower(pack.name)", "%"+strings.ToLower(packName)+"%"))
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
	var filteredSongsCount smdbcore.Count
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

	songs := []smdbcore.Song{}

	err = carta.Map(rows, &songs)
	if err != nil {
		gSugar.Errorln(err)
	}

	c.JSON(200, songs[0])
}
