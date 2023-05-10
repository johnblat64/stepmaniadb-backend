package main

type Chart struct {
	Chartid     string `json:"chartId" `
	Chartname   string `json:"name" `
	StepsType   string `json:"stepsType" `
	Description string `json:"description" `
	Chartstyle  string `json:"chartStyle" `
	Difficulty  string `json:"difficulty" `
	Meter       int    `json:"meter" `
	Credit      string `json:"credit" `
	// Radar Values
	Stream     float64 `json:"stream" `
	Voltage    float64 `json:"voltage" `
	Air        float64 `json:"air" `
	Freeze     float64 `json:"freeze" `
	Chaos      float64 `json:"chaos" `
	TapsCount  int     `json:"tapsCount" db:"taps_count"`
	JumpsCount int     `json:"jumpsCount" db:"jumps_count"`
	HoldsCount int     `json:"holdsCount" db:"holds_count"`
	MinesCount int     `json:"minesCount" db:"mines_count"`
	HandsCount int     `json:"handsCount" db:"hands_count"`
	RollsCount int     `json:"rollsCount" db:"rolls_count"`
	// Gimmicks
	StopsCount   int `json:"stopsCount" db:"stops_count"`
	DelaysCount  int `json:"delaysCount" db:"delays_count"`
	WarpsCount   int `json:"warpsCount" db:"warps_count"`
	ScrollsCount int `json:"scrollsCount" db:"scrolls_count"`
	FakesCount   int `json:"fakesCount" db:"fakes_count"`
	SpeedsCount  int `json:"speedsCount" db:"speeds_count"`
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
	BannerPath     string          `json:"bannerPath" db:"banner_path"`       // looks like "PackName/SongName/banner.png" or "PackName/media/banner.png"
	MusicPath      string          `json:"musicPath" db:"music_path"`         // looks like "PackName/SongName/music.ogg" or "PackName/media/music.ogg"
	SongDirPath    string          `json:"songDirPath" db:"song_dir_path"`    // looks like "PackName/SongName"
	FileExtension  string          `json:"fileExtension" db:"file_extension"` // looks like "sm" or "ssc"
}

type Pack struct {
	PackId       string `json:"packId"   db:"packid"`
	PackName     string `json:"packName" db:"name"`
	DownloadLink string `json:"downloadLink" db:"download_link"`
	BannerPath   string `json:"bannerPath" db:"pack.banner_path"`
	Songs        []Song `json:"songs"`
}

// struct for counting songs
type Count struct {
	Count int `db:"count"`
}
