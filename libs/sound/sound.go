package sound

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"log"
	"os"
	"time"
)

func InitStreamer() beep.StreamSeekCloser {
	f, err := os.Open("/Users/xiusin/projects/src/github.com/xiusin/goqt-translate/misc.mp3")
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	sr := format.SampleRate
	err = speaker.Init(sr, sr.N(time.Second/10))
	return streamer

}
