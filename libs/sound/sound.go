package sound

import (
	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
	"goqt-translate/libs/helper"
	"io/ioutil"
)

func InitStreamer() (*oto.Context, []byte) {
	var file, _ = ioutil.ReadFile(helper.AppDirPath("misc.mp3"))
	dec, data, _ := minimp3.DecodeFull(file)
	ctx, _ := oto.NewContext(dec.SampleRate, dec.Channels, 2, 1024)
	return ctx, data
}
