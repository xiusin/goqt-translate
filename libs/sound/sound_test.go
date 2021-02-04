package sound

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"log"
	"os"
	"testing"
	"time"
)

func TestSound(t *testing.T) {
	var fun = func() {

		f, err := os.Open("/Users/xiusin/projects/src/github.com/xiusin/goqt-translate/misc.mp3")
		if err != nil {
			log.Fatal(err)
		}

		streamer, format, err := mp3.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		defer streamer.Close()

		sr := format.SampleRate

		err = speaker.Init(sr, sr.N(time.Second/10))
		if err != nil {
			t.Error(err)
		}
		for {
			streamer.Seek(0)
			done := make(chan bool)
			//resampled := beep.Resample(4, format.SampleRate, sr, streamer)
			//speaker.Play(resampled)
			speaker.Play(beep.Seq(streamer, beep.Callback(func() {
				done <- true
			})))
			<-done
		}
	}

	fun()
	fmt.Println("播放完成")
}
