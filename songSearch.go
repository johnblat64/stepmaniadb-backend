package main

import "strconv"

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
