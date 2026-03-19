package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/hbollon/go-edlib"

	"github.com/SamW94/get-lyrics/tracks"
)

var (
	featRegex        = regexp.MustCompile(`(?i)\s*(feat\.?|ft\.?|featuring)\s+.*`)
	versionRegex     = regexp.MustCompile(`(?i)\((rmx|remix|live|remaster.*?|demo|edit|version.*?)\)`)
	bracketRegex     = regexp.MustCompile(`\(.*?\)`)
	punctuationRegex = regexp.MustCompile(`[^\w\s]`)
	whitespaceRegex  = regexp.MustCompile(`\s+`)
)

func normaliseSearchStrings(artist, title string) string {

	title = strings.ToLower(title)

	title = featRegex.ReplaceAllString(title, "")
	title = versionRegex.ReplaceAllString(title, "")
	title = bracketRegex.ReplaceAllString(title, "")
	title = punctuationRegex.ReplaceAllString(title, "")
	title = whitespaceRegex.ReplaceAllString(title, " ")

	title = strings.TrimSpace(title)

	return strings.ToLower(artist) + " " + title
}

func slugFromURL(url string) string {
	slug := strings.TrimPrefix(url, "https://genius.com/")
	slug = strings.TrimSuffix(slug, "-lyrics")

	return slug
}

func makeSlug(artist, title string) string {
	slug := normaliseSearchStrings(artist, title)
	slug = strings.ReplaceAll(slug, " ", "-")

	return slug
}

func bestMatchGeniusSearch(query tracks.Track, options []LyricsURL) (LyricsURL, error) {
	querySlug := makeSlug(query.Artist, query.Title)

	bestScore := math.MaxInt
	var bestMatch LyricsURL

	for _, option := range options {
		slug := slugFromURL(option.LyricsURL)

		if slug == querySlug {
			return option, nil
		}

		if strings.Contains(slug, querySlug) || strings.Contains(querySlug, slug) {
			return option, nil
		}

		score := edlib.LevenshteinDistance(querySlug, slug)

		if score < bestScore {
			bestScore = score
			bestMatch = option
		}
	}

	if bestScore > 12 {
		return LyricsURL{}, fmt.Errorf("No acceptable slug match for %v - %v", query.Artist, query.Title)
	}

	return bestMatch, nil
}

func processSearchResponse(geniusSearchResult GeniusSearchResult) ([]LyricsURL, error) {

	if len(geniusSearchResult.Response.Hits) == 0 {
		return nil, fmt.Errorf("No hits found for search term")
	}

	var matches []LyricsURL

	hits := geniusSearchResult.Response.Hits

	if len(hits) > 5 {
		hits = hits[:5]
	}

	for _, hit := range hits {
		if hit.Type != "song" {
			continue
		}

		if !strings.Contains(hit.Result.URL, "-lyrics") {
			continue
		}

		title := strings.TrimSpace(hit.Result.Title)

		matches = append(matches, LyricsURL{
			Artist:    hit.Result.PrimaryArtistNames,
			Title:     title,
			LyricsURL: hit.Result.URL,
		})
	}

	return matches, nil
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

	potentialMatches, err := processSearchResponse(searchResp)
	if err != nil {
		return LyricsURL{}, fmt.Errorf("Error processing search response for %v - %v:\n %v\n", track.Artist, track.Title, err)
	}

	if len(potentialMatches) == 0 {
		return LyricsURL{}, fmt.Errorf("No potential matches found for search term %v - %v", track.Artist, track.Title)
	}

	bestMatch, err := bestMatchGeniusSearch(track, potentialMatches)

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
