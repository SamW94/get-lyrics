package main

import (
	"fmt"
	"os"

	"github.com/SamW94/get-lyrics/tracks"
)

func main() {
	inputArguments := os.Args
	if len(inputArguments) < 2 {
		fmt.Println("No directory provided - please provide one.")
		os.Exit(1)
	}

	directory := inputArguments[1]
	tracks.ReadDirectoryRecursively(directory)
}
