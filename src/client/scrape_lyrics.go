package client

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/SamW94/get-lyrics/tracks"
	"github.com/gocolly/colly"
)

func (c *Client) ScrapeLyrics(track tracks.Track) (lyrics string, err error) {
	collyCollector := colly.NewCollector(
		colly.AllowedDomains("www.genius.com", "genius.com"),
	)

	collyCollector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"

	var lyricsParagraphs []string

	collyCollector.OnHTML(`div[data-lyrics-container="true"]`, func(e *colly.HTMLElement) {

		e.DOM.Find("br").ReplaceWithHtml("\n")
		text := strings.TrimSpace(e.DOM.Text())
		if text != "" {
			lyricsParagraphs = append(lyricsParagraphs, text)
		}
	})

	collyCollector.OnError(func(r *colly.Response, err error) {
		log.Printf("Scraping failed for %s: %v\n", track.LyricsURL, err)
	})

	err = collyCollector.Visit(track.LyricsURL)
	if err != nil {
		return "", fmt.Errorf("Error scraping lyrics for %v - %v:\n%v", track.Artist, track.Title, err)
	}

	lyrics = strings.Join(lyricsParagraphs, "\n")
	if idx := strings.Index(lyrics, "["); idx != -1 {
		lyrics = lyrics[idx:]
	}

	return lyrics, nil
}

func (c *Client) ScrapeLyricsConcurrent(trackList []tracks.Track) (successful []tracks.Track, failed []string) {
	var waitGroup sync.WaitGroup
	var mutex sync.Mutex
	jobs := make(chan tracks.Track)
	workerCount := 5
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	var tracksSuccessful []tracks.Track
	var trackArtistTitleFailed []string
	for worker := 0; worker < workerCount; worker++ {
		waitGroup.Add(1)

		go func(workerID int) {

			defer waitGroup.Done()

			for track := range jobs {
				<-ticker.C

				lyrics, err := c.ScrapeLyrics(track)
				if err != nil {
					log.Printf("Error scraping lyrics:\n %v\n", err)

					mutex.Lock()
					trackArtistTitleFailed = append(trackArtistTitleFailed, fmt.Sprintf("%v - %v", track.Artist, track.Title))
					mutex.Unlock()
					continue
				}
				track.Lyrics = lyrics

				mutex.Lock()
				tracksSuccessful = append(tracksSuccessful, track)
				mutex.Unlock()
			}
		}(worker)
	}

	for _, track := range trackList {
		jobs <- track
	}
	close(jobs)
	waitGroup.Wait()

	return tracksSuccessful, trackArtistTitleFailed
}
