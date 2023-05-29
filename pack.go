package main

import (
	"github.com/gin-gonic/gin"
	sqlBuilder "github.com/huandu/go-sqlbuilder"
	"github.com/jackskj/carta"
	smdbcore "github.com/stepmaniadb/stepmaniadb-core"
)

func getPackById(c *gin.Context) {
	packid := c.Param("packid")

	sqlPackSelectBuilder := sqlBuilder.NewSelectBuilder()
	sqlPackSelectBuilder.Select("song.songid, song.title, song.artist,  song.banner_path, song.music_path, song.song_dir_path, song_bpm.song_bpm, song_time_signature.time_signature_numerator, song_time_signature.time_signature_denominator, chart.meter, pack.packid, pack.name, pack.download_link, pack.pack_banner_path")
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

	packs := []smdbcore.Pack{}

	err = carta.Map(rows, &packs)
	if err != nil {
		gSugar.Errorln(err)
	}

	c.JSON(200, packs[0])
}
