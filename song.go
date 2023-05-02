package main

import "strconv"

type Chart struct {
	Chartid      string `json:"chartId" `
	Chartname    string `json:"name" `
	StepsType    string `json:"stepsType" `
	Description  string `json:"description" `
	Chartstyle   string `json:"chartStyle" `
	Difficulty   string `json:"difficulty" `
	Meter        int    `json:"meter" `
	Credit       string `json:"credit" `
	StopsCount   int    `json:"stopsCount" db:"stops_count"`
	DelaysCount  int    `json:"delaysCount" db:"delays_count"`
	WarpsCount   int    `json:"warpsCount" db:"warps_count"`
	ScrollsCount int    `json:"scrollsCount" db:"scrolls_count"`
	FakesCount   int    `json:"fakesCount" db:"fakes_count"`
	SpeedsCount  int    `json:"speedsCount" db:"speeds_count"`
}

type Bpm struct {
	Value float32 `json:"value" db:"song_bpm"`
}

type TimeSignature struct {
	Numerator   int `json:"numerator" db:"time_signature_numerator"`
	Denominator int `json:"denominator" db:"time_signature_denominator"`
}

type Song struct {
	SongId         string          `json:"songId"         `
	Title          string          `json:"title"      `
	Artist         string          `json:"artist"     `
	Bpms           []Bpm           `json:"bpms"           `
	TimeSignatures []TimeSignature `json:"timeSignatures" `
	Charts         []Chart         `json:"charts"`
	PackId         string          `json:"packId"         db:"packid"`
	PackName       string          `json:"packName"       db:"name"`
}

type SongNugget struct {
	SongId         string          `json:"songId"   db:"songid"`
	Title          string          `json:"title"    db:"title"`
	Artist         string          `json:"artist"   db:"artist"`
	Bpms           []Bpm           `json:"bpms"           `
	TimeSignatures []TimeSignature `json:"timeSignatures" `
}

type Pack struct {
	PackId       string `json:"packId"   db:"packid"`
	PackName     string `json:"packName" db:"name"`
	DownloadLink string `json:"downloadLink" db:"download_link"`
	Songs        []Song `json:"songs"`
}

type SongResultsResponse struct {
	Page            int    `json:"pageNum"`
	PageSize        int    `json:"pageSize"`
	PageCount       int    `json:"pageCount"`
	TotalSongsCount int    `json:"totalSongsCount" db:"total_songs_count"`
	Songs           []Song `json:"songs"`
}

type SongSearchParameters struct {
	Title                    string
	Artist                   string
	Credit                   string
	Pack                     string
	TimeSignatureNumerator   int
	TimeSignatureDenominator int
	BpmMin                   int
	BpmMax                   int
	MeterMin                 int
	MeterMax                 int
	StepsType                string
}

// for HTML page
type SongsResultsModel struct {
	Page             int    `json:"pageNum"`
	PageSize         int    `json:"pageSize"`
	PageCount        int    `json:"pageCount"`
	TotalSongsCount  int    `json:"totalSongsCount"`
	Songs            []Song `json:"songs"`
	StepsTypeOptions []string
	SearchParameters SongSearchParameters
}

type SongModel struct {
	Song Song `json:"song"`
}

func (params SongSearchParameters) AsQueryString() string {
	return "?title=" + params.Title + "&artist=" + params.Artist + "&credit=" + params.Credit + "&pack=" + params.Pack + "&stepstype=" + params.StepsType + "&timeSignatureNumerator=" + strconv.Itoa(params.TimeSignatureNumerator) + "&timeSignatureDenominator=" + strconv.Itoa(params.TimeSignatureDenominator) + "&bpmMin=" + strconv.Itoa(params.BpmMin) + "&bpmMax=" + strconv.Itoa(params.BpmMax) + "&meterMin=" + strconv.Itoa(params.MeterMin) + "&meterMax=" + strconv.Itoa(params.MeterMax)
}

func (resultsModel SongsResultsModel) NextPage() int {
	return resultsModel.Page + 1
}

func (resultsModel SongsResultsModel) HasNextPage() bool {
	return resultsModel.Page < resultsModel.PageCount
}

func (resultsModel SongsResultsModel) PreviousPage() int {
	return resultsModel.Page - 1
}

// struct for counting songs
type Count struct {
	Count int `db:"count"`
}
