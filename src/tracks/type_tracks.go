package tracks

type FilePath struct {
	Path *string
}

type Track struct {
	Path      string
	Artist    string
	Title     string
	LyricsURL string
	Lyrics    string
}
