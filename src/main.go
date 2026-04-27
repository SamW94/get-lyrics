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
	tracksList, failedMp3s, err := tracks.GetTracks(directory)

	webClient := client.NewClient(time.Second * 5)

	tracksSearchSuccessful, tracksSearchFailed := webClient.GeniusSearchConcurrent(tracksList)
	tracksLyricsSuccessful, trackLyricsFailed := webClient.ScrapeLyricsConcurrent(tracksSearchSuccessful)

	for _, t := range tracksLyricsSuccessful {
		err = tracks.TagMp3(t)
		if err != nil {
			fmt.Printf("Failed to tag MP3 at %v: %v", t.Path, err)
		}
	}

	for _, failedMp3 := range failedMp3s {
		fmt.Printf("Failed to read the artist or title tag from mp3 file %v - is there a problem with the tag?", failedMp3)
	}

	for _, trackArtistTitle := range tracksSearchFailed {
		fmt.Printf("Failed to find any good hits when searching for %v\n", trackArtistTitle)
	}

	for _, trackArtistTitle := range trackLyricsFailed {
		fmt.Printf("Failed to scrape lyrics for %v\n", trackArtistTitle)
	}

	fmt.Println("Completed!")
}
