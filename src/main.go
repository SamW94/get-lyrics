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
	tracksList, err := tracks.ListTracksForMP3s(directory)
	tracksWithoutLyricsList := tracks.FindTracksWithoutLyrics(tracksList)

	webClient := client.NewClient(time.Second * 5)

	lyricsURLs := webClient.GeniusSearchConcurrent(tracksWithoutLyricsList)
	for _, lyricsURL := range lyricsURLs {
		fmt.Printf("%v\n", lyricsURL)
	}
}
