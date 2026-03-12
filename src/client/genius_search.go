package client

import (
	"encoding/json"
	"fmt"
	"io"
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

func processSearchResponse(geniusSearchResult GeniusSearchResult) (potentialMatches []LyricsURL, err error) {

	if len(geniusSearchResult.Response.Hits) == 0 {
		return potentialMatches, fmt.Errorf("No hits found for search term")
	}

	for _, hit := range geniusSearchResult.Response.Hits {
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
	request, err := constructHTTPRequest(url, searchTerm)
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

	potentialMatches, err := processSearchResponse(searchResp)
	if err != nil {
		return LyricsURL{}, fmt.Errorf("Error processing search response: %v\n", err)
	}

	bestMatch, err := bestMatchGeniusSearch(track.Title, potentialMatches)

	if err != nil {
		return LyricsURL{}, fmt.Errorf("No matches found for %v - %v: %v\n", track.Artist, track.Title, err)
	}

	return bestMatch, nil
}

func (c *Client) GeniusSearchConcurrent(trackList []tracks.Track) []LyricsURL {
	var waitGroup sync.WaitGroup
	jobs := make(chan tracks.Track)
	workerCount := 5
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	var lyricsURLs []LyricsURL
	for worker := 0; worker < workerCount; worker++ {
		waitGroup.Add(1)

		go func(workerID int) {

			defer waitGroup.Done()

			for track := range jobs {
				<-ticker.C

				lyricsURL, err := c.GeniusSearch(track)
				if err != nil {
					fmt.Printf("Error searching for lyrics:\n %v\n", err)
					continue
				}
				lyricsURLs = append(lyricsURLs, lyricsURL)

			}
		}(worker)
	}

	for _, track := range trackList {
		jobs <- track
	}
	close(jobs)
	waitGroup.Wait()

	return lyricsURLs
}
