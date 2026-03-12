package main

import (
	"fmt"
	"os"
	"time"

	"github.com/SamW94/get-lyrics/client"
	"github.com/SamW94/get-lyrics/tracks"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Printf("Failed to load environment vars from .env file: %v\n", err)
		os.Exit(1)
	}

	inputArguments := os.Args
	if len(inputArguments) < 2 {
		fmt.Println("No directory provided - please provide one.")
		os.Exit(1)
	}

	directory := inputArguments[1]
	tracksList, err := tracks.GetTracks(directory)

	webClient := client.NewClient(time.Second * 5)

	tracksSearchSuccessful, tracksSearchFailed := webClient.GeniusSearchConcurrent(tracksList)
	tracksLyricsSuccessful, trackLyricsFailed := webClient.ScrapeLyricsConcurrent(tracksSearchSuccessful)

	for _, track := range tracksLyricsSuccessful {
		fmt.Printf("%v\n", track)
	}

	for _, trackArtistTitle := range tracksSearchFailed {
		fmt.Printf("Failed to find any good hits when searching for %v\n", trackArtistTitle)
	}

	for _, trackArtistTitle := range trackLyricsFailed {
		fmt.Printf("Failed to scrape lyrics for %v\n", trackArtistTitle)
	}
}
