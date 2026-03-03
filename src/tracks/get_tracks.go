package tracks

import (
	"fmt"
	"html"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"go.senan.xyz/taglib"
)

func createTrackObject(mp3 string) (Track, error) {

	tags, err := taglib.ReadTags(mp3)
	if err != nil {
		fmt.Printf("Error reading MP3 tags from file %v: %v", mp3, err)
		return Track{}, err
	}

	artist := tags[taglib.Artist][0]
	title := tags[taglib.Title][0]
	lyrics := tags[taglib.Lyrics][0]

	track := Track{
		Artist: artist,
		Title:  title,
		Lyrics: lyrics,
	}

	return track, nil

}

func createTrackObjectsConcurrently(mp3s []string) ([]Track, error) {
	var waitGroup sync.WaitGroup
	jobs := make(chan string, len(mp3s))
	trackObjects := make(chan Track, len(mp3s))

	cpuWorkers := runtime.NumCPU()

	for range cpuWorkers {
		waitGroup.Go(func() {
			for mp3 := range jobs {
				track, err := createTrackObject(mp3)
				if err != nil {
					fmt.Printf("Error creating track object for %v: %v", mp3, err)
				}
				trackObjects <- track
			}
		})
	}

	for _, mp3 := range mp3s {
		jobs <- mp3
	}
	close(jobs)

	go func() {
		waitGroup.Wait()
		close(trackObjects)
	}()

	var tracks []Track
	for trackObject := range trackObjects {
		tracks = append(tracks, trackObject)
	}

	return tracks, nil
}

func getMp3s(directory string) ([]string, error) {
	var mp3s []string

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".mp3") && !strings.Contains(path, "Various Artists") {
			mp3s = append(mp3s, path)
		}
		return nil
	})

	return mp3s, err
}

func ReadDirectoryRecursively(directory string) {
	fmt.Printf("Finding MP3 files in the %v directory and its subdirectories, please wait...\n", directory)
	mp3s, err := getMp3s(directory)
	if err != nil {
		fmt.Printf("Error retrieving MP3s from directory: %v", err)
	}
	if len(mp3s) == 0 {
		fmt.Println("No MP3 files found after searching directory recursively.")
	}

	fmt.Printf("%d MP3 files found in directory %v\n", len(mp3s), directory)
	for _, mp3 := range mp3s {
		fmt.Printf("%v\n", mp3)
	}

	tracks, err := createTrackObjectsConcurrently(mp3s)
	for _, track := range tracks {

		if track.Title == "Scarred" {
			lyrics, err := getLyrics("https://www.last.fm/music/Dream+Theater/_/Scarred/+lyrics")
			if err != nil {
				fmt.Printf("%v", err)
			}
			fmt.Println(lyrics)
		}
	}

}

func getLyrics(url string) (string, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.last.fm", "last.fm"),
	)

	var paragraphs []string

	c.OnHTML("span.lyrics-body p.lyrics-paragraph", func(e *colly.HTMLElement) {
		var builder strings.Builder

		e.DOM.Contents().Each(func(i int, s *goquery.Selection) {
			if goquery.NodeName(s) == "br" {
				builder.WriteString("\n")
			} else {
				builder.WriteString(s.Text())
			}
		})

		text := strings.TrimSpace(builder.String())
		text = html.UnescapeString(text)

		if text != "" {
			paragraphs = append(paragraphs, text)
		}
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	lyrics := strings.Join(paragraphs, "\n\n")
	return lyrics, nil
}
