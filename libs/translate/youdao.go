package translate

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/hajimehoshi/oto"
	request2 "github.com/mozillazg/request"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/tosone/minimp3"
	"github.com/xiusin/logger"
	"goqt-translate/libs/request"
	"net/url"
	"time"
	"unicode/utf8"
)

const youdaoApiUrl = "https://openapi.youdao.com/api"

type YouDaoTranslate struct {
	AppId    string
	Secret   string
	soundCtx *oto.Context
}

func (b *YouDaoTranslate) Translate(from, to, q string) (string, error) {
	uuidGen, _ := uuid.NewV4()
	salt := uuidGen.String()
	curTime := fmt.Sprint(time.Now().Unix())

	sign := b.Sign(salt, curTime, q)
	req := request.RequestPool.Get()
	defer request.RequestPool.Put(req)
	r := req.(*request2.Request)
	u := url.Values{
		"appKey":   []string{b.AppId},
		"q":        []string{q},
		"from":     []string{from},
		"to":       []string{to},
		"salt":     []string{salt},
		"sign":     []string{sign},
		"signType": []string{"v3"},
		"curtime":  []string{curTime},
	}
	resp, err := r.Get(youdaoApiUrl + "?" + u.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := resp.Json()
	if err != nil {
		return "", err
	}
	errorMsg, ok := data.CheckGet("error_msg")
	if ok == true {
		msg, _ := errorMsg.String()
		return "", errors.New(msg)
	}
	d, err := data.Get("translation").GetIndex(0).String()
	return d, err
}

func TranslateYoudao(from, to, q string) (string, error) {
	return youDaoTranslateStd.Translate(from, to, q)
}

func ToVoice(str, lang string) {
	youDaoTranslateStd.ToVoice(str, lang)
}

func (b *YouDaoTranslate) Sign(salt, curTime, q string) string {
	l := utf8.RuneCountInString(q)
	input := ""
	if l <= 20 {
		input = q
	} else {
		run := []rune(q)
		input = fmt.Sprintf("%s%d%s", string(run[:10]), l, string(run[l-10:]))
	}
	has := sha256.Sum256([]byte(b.AppId + input + salt + curTime + b.Secret))
	return fmt.Sprintf("%x", has)
}

func (b *YouDaoTranslate) ToVoice(str, lang string) error {
	if lang == "en" {
		lang = "eng"
	} else {
		lang = "zh"
	}
	var url = fmt.Sprintf("http://dict.youdao.com/dictvoice?audio=%s&le=%s", url.QueryEscape(str), lang)
	req := request.RequestPool.Get()
	defer request.RequestPool.Put(req)
	r := req.(*request2.Request)
	resp, err := r.Post(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := resp.Content()
	if err != nil {
		return err
	}
	dec, decodeData, _ := minimp3.DecodeFull(data)
	if dec.SampleRate != 0 {
		data = decodeData
	} else {
		dec.SampleRate, dec.Channels = 24000, 1
	}
	soundCtx, err := oto.NewContext(dec.SampleRate, dec.Channels, 2, 1024)
	if err != nil {
		logger.Error(err)
	} else {
		defer soundCtx.Close()
		p := soundCtx.NewPlayer()
		p.Write(data)
		p.Close()
	}

	return nil
}

var youDaoTranslateStd *YouDaoTranslate

func init() {
	youDaoTranslateStd = &YouDaoTranslate{
		AppId:  "55ef1f66b723fcf5",
		Secret: "38D3HGvTeMqlQ9M53KVh4YmcsZp2hgk0",
	}
}
