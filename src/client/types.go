package client

import (
	"net/http"
)

type Client struct {
	httpClient http.Client
}

type GeniusSearchResult struct {
	Meta struct {
		Status int `json:"status"`
	} `json:"meta"`
	Response struct {
		Hits []struct {
			Highlights []interface{} `json:"highlights"`
			Index      string        `json:"index"`
			Type       string        `json:"type"`
			Result     struct {
				AnnotationCount         int    `json:"annotation_count"`
				APIPath                 string `json:"api_path"`
				ArtistNames             string `json:"artist_names"`
				FullTitle               string `json:"full_title"`
				HeaderImageThumbnailURL string `json:"header_image_thumbnail_url"`
				HeaderImageURL          string `json:"header_image_url"`
				ID                      int    `json:"id"`
				LyricsOwnerID           int    `json:"lyrics_owner_id"`
				LyricsState             string `json:"lyrics_state"`
				Path                    string `json:"path"`
				PrimaryArtistNames      string `json:"primary_artist_names"`
				PyongsCount             int    `json:"pyongs_count"`
				RelationshipsIndexURL   string `json:"relationships_index_url"`
				ReleaseDateComponents   struct {
					Year  int `json:"year"`
					Month int `json:"month"`
					Day   int `json:"day"`
				} `json:"release_date_components"`
				ReleaseDateForDisplay                     string `json:"release_date_for_display"`
				ReleaseDateWithAbbreviatedMonthForDisplay string `json:"release_date_with_abbreviated_month_for_display"`
				SongArtImageThumbnailURL                  string `json:"song_art_image_thumbnail_url"`
				SongArtImageURL                           string `json:"song_art_image_url"`
				Stats                                     struct {
					UnreviewedAnnotations int  `json:"unreviewed_annotations"`
					Hot                   bool `json:"hot"`
					Pageviews             int  `json:"pageviews"`
				} `json:"stats"`
				Title             string        `json:"title"`
				TitleWithFeatured string        `json:"title_with_featured"`
				URL               string        `json:"url"`
				FeaturedArtists   []interface{} `json:"featured_artists"`
				PrimaryArtist     struct {
					APIPath        string `json:"api_path"`
					HeaderImageURL string `json:"header_image_url"`
					ID             int    `json:"id"`
					ImageURL       string `json:"image_url"`
					IsMemeVerified bool   `json:"is_meme_verified"`
					IsVerified     bool   `json:"is_verified"`
					Name           string `json:"name"`
					URL            string `json:"url"`
				} `json:"primary_artist"`
				PrimaryArtists []struct {
					APIPath        string `json:"api_path"`
					HeaderImageURL string `json:"header_image_url"`
					ID             int    `json:"id"`
					ImageURL       string `json:"image_url"`
					IsMemeVerified bool   `json:"is_meme_verified"`
					IsVerified     bool   `json:"is_verified"`
					Name           string `json:"name"`
					URL            string `json:"url"`
				} `json:"primary_artists"`
			} `json:"result"`
		} `json:"hits"`
	} `json:"response"`
}

type GeniusSearchResultHit struct {
	Highlights []interface{} `json:"highlights"`
	Index      string        `json:"index"`
	Type       string        `json:"type"`
	Result     struct {
		AnnotationCount         int    `json:"annotation_count"`
		APIPath                 string `json:"api_path"`
		ArtistNames             string `json:"artist_names"`
		FullTitle               string `json:"full_title"`
		HeaderImageThumbnailURL string `json:"header_image_thumbnail_url"`
		HeaderImageURL          string `json:"header_image_url"`
		ID                      int    `json:"id"`
		LyricsOwnerID           int    `json:"lyrics_owner_id"`
		LyricsState             string `json:"lyrics_state"`
		Path                    string `json:"path"`
		PrimaryArtistNames      string `json:"primary_artist_names"`
		PyongsCount             int    `json:"pyongs_count"`
		RelationshipsIndexURL   string `json:"relationships_index_url"`
		ReleaseDateComponents   struct {
			Year  int `json:"year"`
			Month int `json:"month"`
			Day   int `json:"day"`
		} `json:"release_date_components"`
		ReleaseDateForDisplay                     string `json:"release_date_for_display"`
		ReleaseDateWithAbbreviatedMonthForDisplay string `json:"release_date_with_abbreviated_month_for_display"`
		SongArtImageThumbnailURL                  string `json:"song_art_image_thumbnail_url"`
		SongArtImageURL                           string `json:"song_art_image_url"`
		Stats                                     struct {
			UnreviewedAnnotations int  `json:"unreviewed_annotations"`
			Hot                   bool `json:"hot"`
			Pageviews             int  `json:"pageviews"`
		} `json:"stats"`
		Title             string        `json:"title"`
		TitleWithFeatured string        `json:"title_with_featured"`
		URL               string        `json:"url"`
		FeaturedArtists   []interface{} `json:"featured_artists"`
		PrimaryArtist     struct {
			APIPath        string `json:"api_path"`
			HeaderImageURL string `json:"header_image_url"`
			ID             int    `json:"id"`
			ImageURL       string `json:"image_url"`
			IsMemeVerified bool   `json:"is_meme_verified"`
			IsVerified     bool   `json:"is_verified"`
			Name           string `json:"name"`
			URL            string `json:"url"`
		} `json:"primary_artist"`
		PrimaryArtists []struct {
			APIPath        string `json:"api_path"`
			HeaderImageURL string `json:"header_image_url"`
			ID             int    `json:"id"`
			ImageURL       string `json:"image_url"`
			IsMemeVerified bool   `json:"is_meme_verified"`
			IsVerified     bool   `json:"is_verified"`
			Name           string `json:"name"`
			URL            string `json:"url"`
		} `json:"primary_artists"`
	} `json:"result"`
}

type LyricsURL struct {
	Artist    string
	Title     string
	LyricsURL string
}
