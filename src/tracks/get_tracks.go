package tracks

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/bogem/id3v2"
)

func createTrackObject(mp3 string) (Track, error) {

	tag, err := id3v2.Open(mp3, id3v2.Options{Parse: true})
	if err != nil {
		return Track{}, fmt.Errorf("Error opening file at %v:\n %v", mp3, err)
	}

	defer tag.Close()

	artist := tag.Artist()
	title := tag.Title()

	if artist == "" {
		return Track{}, fmt.Errorf("Error obtaining artist from tags for file %v - is there a problem with the tag?\n", mp3)
	}

	if title == "" {
		return Track{}, fmt.Errorf("Error obtaining title from tags for file %v - is there a problem with the tag?\n", mp3)
	}

	track := Track{
		Path:   mp3,
		Artist: artist,
		Title:  title,
		Lyrics: "",
	}

	return track, nil

}

func createTrackObjectsConcurrently(mp3s []string) (successful []Track, failed []string) {
	var waitGroup sync.WaitGroup
	var mutex sync.Mutex
	jobs := make(chan string, len(mp3s))
	trackObjects := make(chan Track, len(mp3s))

	cpuWorkers := runtime.NumCPU()

	var tracks []Track
	var failedMp3s []string

	for i := 0; i < cpuWorkers; i++ {
		waitGroup.Go(func() {
			for mp3 := range jobs {
				track, err := createTrackObject(mp3)
				if err != nil {
					log.Printf("Error creating track object for %v\n: %v\n", mp3, err)
					mutex.Lock()
					failedMp3s = append(failedMp3s, mp3)
					mutex.Unlock()
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

	for trackObject := range trackObjects {
		mutex.Lock()
		tracks = append(tracks, trackObject)
		mutex.Unlock()
	}

	return tracks, failedMp3s
}

func getMp3s(directory string) ([]string, error) {
	var mp3s []string

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".mp3") {
			mp3s = append(mp3s, path)
		}
		return nil
	})

	return mp3s, err
}

func GetTracks(directory string) ([]Track, []string, error) {
	fmt.Printf("Finding MP3 files in the %v directory and its subdirectories, please wait...\n", directory)
	mp3Paths, err := getMp3s(directory)
	if err != nil {
		return nil, nil, fmt.Errorf("Error retrieving MP3s from directory\n: %v\n", err)
	}
	if len(mp3Paths) == 0 {
		return nil, nil, fmt.Errorf("No MP3 files found after searching directory recursively.")
	}

	fmt.Printf("%d MP3 files found in directory %v\n", len(mp3Paths), directory)

	tracks, failed := createTrackObjectsConcurrently(mp3Paths)
	return tracks, failed, nil
}
