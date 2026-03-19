package tracks

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"go.senan.xyz/taglib"
)

func findTracksWithoutLyrics(tracks []Track) (tracksNoLyrics []Track) {
	for _, track := range tracks {
		if track.Lyrics == "" {
			tracksNoLyrics = append(tracksNoLyrics, track)
		}
	}

	return tracksNoLyrics
}

func createTrackObject(mp3 string) (Track, error) {

	tags, err := taglib.ReadTags(mp3)
	if err != nil {
		return Track{}, fmt.Errorf("Error reading MP3 tags from file %v: %v", mp3, err)
	}

	var lyrics string
	var artist string
	var title string

	if len(tags[taglib.Artist]) != 0 {
		artist = tags[taglib.Artist][0]
	} else {
		return Track{}, fmt.Errorf("Error obtaining artist from tags for file %v - is there a problem with the tag?\n", mp3)
	}

	if len(tags[taglib.Title]) != 0 {
		title = tags[taglib.Title][0]
	} else {
		return Track{}, fmt.Errorf("Error obtaining title from tags for file %v - is there a problem with the tag?\n", mp3)
	}

	if len(tags[taglib.Lyrics]) != 0 {
		lyrics = tags[taglib.Lyrics][0]
	}
	lyrics = ""

	track := Track{
		Path:   mp3,
		Artist: artist,
		Title:  title,
		Lyrics: lyrics,
	}
	fmt.Printf("%v\n", track)

	return track, nil

}

func createTrackObjectsConcurrently(mp3s []string) (successful []Track, failed []string, err error) {
	var waitGroup sync.WaitGroup
	var mutex sync.Mutex
	jobs := make(chan string, len(mp3s))
	trackObjects := make(chan Track, len(mp3s))

	cpuWorkers := runtime.NumCPU()

	var tracks []Track
	var failedMp3s []string

	for range cpuWorkers {
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

	return tracks, failedMp3s, nil
}

func getMp3s(directory string) ([]string, error) {
	var mp3s []string

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".mp3") && !strings.Contains(path, "Essential Mix") {
			mp3s = append(mp3s, path)
		}
		return nil
	})

	return mp3s, err
}

func GetTracks(directory string) ([]Track, []string, error) {
	fmt.Printf("Finding MP3 files in the %v directory and its subdirectories, please wait...\n", directory)
	mp3s, err := getMp3s(directory)
	if err != nil {
		return nil, nil, fmt.Errorf("Error retrieving MP3s from directory\n: %v\n", err)
	}
	if len(mp3s) == 0 {
		return nil, nil, fmt.Errorf("No MP3 files found after searching directory recursively.")
	}

	fmt.Printf("%d MP3 files found in directory %v\n", len(mp3s), directory)
	for _, mp3 := range mp3s {
		fmt.Printf("%v\n", mp3)
	}

	tracks, failed, err := createTrackObjectsConcurrently(mp3s)
	if err != nil {
		return nil, nil, fmt.Errorf("Error creating track objects from list of MP3s\n: %v\n", err)
	}

	tracks = findTracksWithoutLyrics(tracks)

	return tracks, failed, nil
}
