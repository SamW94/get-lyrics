package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/hbollon/go-edlib"

	"github.com/SamW94/get-lyrics/tracks"
)

func bestMatchGeniusSearch(query string, options []LyricsURL) (LyricsURL, error) {
	titleOptionsMap := make(map[string]LyricsURL)
	var titleOptionSlice []string

	for _, option := range options {
		titleOptionsMap[option.Title] = option
	}

	for title := range titleOptionsMap {
		titleOptionSlice = append(titleOptionSlice, title)
	}

	bestMatchString, err := edlib.FuzzySearch(query, titleOptionSlice, edlib.Levenshtein)
	if err != nil {
		return LyricsURL{}, fmt.Errorf("Fuzzy search for best match failed: %v", err)
	}

	return titleOptionsMap[bestMatchString], nil

}

func processSearchResponse(geniusSearchResult GeniusSearchResult, track tracks.Track) (potentialMatches []LyricsURL, err error) {

	if len(geniusSearchResult.Response.Hits) == 0 {
		return nil, fmt.Errorf("No hits found for search term")
	}

	for _, hit := range geniusSearchResult.Response.Hits {
		if hit.Type != "song" {
			continue
		}

		if !strings.Contains(hit.Result.PrimaryArtistNames, track.Artist) {
			continue
		}

		lyricsURL := LyricsURL{
			Artist:    hit.Result.PrimaryArtistNames,
			Title:     hit.Result.TitleWithFeatured,
			LyricsURL: hit.Result.URL,
		}
		potentialMatches = append(potentialMatches, lyricsURL)
	}
	return potentialMatches, nil
}

func (c *Client) GeniusSearch(track tracks.Track) (lyricsURL LyricsURL, err error) {
	searchTerm := fmt.Sprintf("%v %v", track.Artist, track.Title)
	url := baseURL + "/search"
	request, err := constructHTTPRequestSearch(url, searchTerm)
	if err != nil {
		return LyricsURL{}, fmt.Errorf("Error constructing HTTP request\n: %v\n", err)
	}

	response, err := c.makeHTTPRequest(*request)
	if err != nil {
		return LyricsURL{}, fmt.Errorf("Error making HTTP request\n: %v\n", err)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return LyricsURL{}, fmt.Errorf("Error reading response body from %v\n: %v\n", url, err)
	}

	defer response.Body.Close()

	searchResp := GeniusSearchResult{}
	err = json.Unmarshal(data, &searchResp)
	if err != nil {
		return LyricsURL{}, fmt.Errorf("Error unmarshalling search results: %v\n", err)
	}

	potentialMatches, err := processSearchResponse(searchResp, track)
	if err != nil {
		return LyricsURL{}, fmt.Errorf("Error processing search response for %v - %v:\n %v\n", track.Artist, track.Title, err)
	}

	var filteredPotentialMatches []LyricsURL
	for _, potentialMatch := range potentialMatches {
		if potentialMatch.Artist == track.Artist {
			filteredPotentialMatches = append(filteredPotentialMatches, potentialMatch)
		}
	}

	if len(filteredPotentialMatches) == 0 {
		return LyricsURL{}, fmt.Errorf("No potential matches found for search term %v - %v", track.Artist, track.Title)
	}

	bestMatch, err := bestMatchGeniusSearch(track.Title, filteredPotentialMatches)

	if err != nil {
		return LyricsURL{}, fmt.Errorf("No matches found for %v - %v: %v\n", track.Artist, track.Title, err)
	}

	return bestMatch, nil
}

func (c *Client) GeniusSearchConcurrent(trackList []tracks.Track) (successful []tracks.Track, failed []string) {
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

				lyricsURL, err := c.GeniusSearch(track)
				if err != nil {
					log.Printf("Error searching for lyrics:\n %v\n", err)
					mutex.Lock()
					trackArtistTitleFailed = append(trackArtistTitleFailed, fmt.Sprintf("%v - %v", track.Artist, track.Title))
					mutex.Unlock()
					continue
				}
				track.LyricsURL = lyricsURL.LyricsURL

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
