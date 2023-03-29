package main

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
