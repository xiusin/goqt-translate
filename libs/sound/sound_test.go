package sound

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
	"io/ioutil"
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

func TestMiniMp3(t *testing.T) {
	var file, _ = ioutil.ReadFile("/Users/xiusin/projects/src/github.com/xiusin/goqt-translate/misc.mp3")
	dec, data, _ := minimp3.DecodeFull(file)
	ctx, _ := oto.NewContext(dec.SampleRate, dec.Channels, 2, 1024)
	p := ctx.NewPlayer()
	go func() {
		for {
			p.Write(data)
			time.Sleep(20 * time.Millisecond)
		}
	}()

	go func() {
		for {
			p.Write(data)
			time.Sleep(30 * time.Millisecond)
		}
	}()

	select {}
}
