package tracks

func FindTracksWithoutLyrics(tracks []Track) (tracksNoLyrics []Track) {
	for _, track := range tracks {
		if track.Lyrics == "" {
			tracksNoLyrics = append(tracksNoLyrics, track)
		}
	}

	return tracksNoLyrics
}
