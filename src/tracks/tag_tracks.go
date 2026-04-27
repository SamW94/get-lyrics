package tracks

import (
	"fmt"

	"github.com/bogem/id3v2"
)

func TagMp3(t Track) error {
	tag, err := id3v2.Open(t.Path, id3v2.Options{Parse: true})
	if err != nil {
		return fmt.Errorf("Error opening MP3 file at path %v:\n %v", t.Path, err)
	}
	defer tag.Close()

	tag.DeleteFrames("USLT")
	uslf := id3v2.UnsynchronisedLyricsFrame{
		Encoding:          id3v2.EncodingUTF8,
		Language:          "eng",
		ContentDescriptor: "Lyrics",
		Lyrics:            t.Lyrics,
	}

	tag.AddUnsynchronisedLyricsFrame(uslf)

	err = tag.Save()
	if err != nil {
		return fmt.Errorf("Error writing lyrics tag to file %v:\n %v", t.Path, err)
	}

	return nil
}
