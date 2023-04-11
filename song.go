package main

type Chart struct {
	Chartid     string `json:"chartId" `
	Chartname   string `json:"name" `
	Stepstype   string `json:"stepsType" `
	Description string `json:"description" `
	Chartstyle  string `json:"chartStyle" `
	Difficulty  string `json:"difficulty" `
	Meter       int    `json:"meter" `
	Credit      string `json:"credit" `
	// StopsCount   int    `json:"stopsCount" db:"chart.stops_count"`
	// DelaysCount  int    `json:"delaysCount" db:"chart.delays_count"`
	// WarpsCount   int    `json:"warpsCount" db:"chart.warps_count"`
	// ScrollsCount int    `json:"scrollsCount" db:"chart.scrolls_count"`
	// FakesCount   int    `json:"fakesCount" db:"chart.fakes_count"`
	// SpeedsCount  int    `json:"speedsCount" db:"chart.speeds_count"`
}

type Bpm struct {
	Value float32 `db:"song_bpm"`
}

type TimeSignature struct {
	Numerator   int `db:"time_signature_numerator"`
	Denominator int `db:"time_signature_denominator"`
}

type Song struct {
	Songid         string          `json:"songId"         `
	Title          string          `json:"songTitle"      `
	Artist         string          `json:"songArtist"     `
	Bpms           []Bpm           `json:"bpms"           `
	Timesignatures []TimeSignature `json:"timeSignatures" `
	Charts         []Chart         `json:"charts"`

	// 	PackId         string          `json:"packId"         db:"pack.packid"`
	//	PackName       string          `json:"packName"       db:"pack.name"`

}
