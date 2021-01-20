package translate

import (
	"crypto/sha256"
	"errors"
	"fmt"
	request2 "github.com/mozillazg/request"
	uuid "github.com/nu7hatch/gouuid"
	"goqt-translate/libs/request"
	"net/url"
	"time"
	"unicode/utf8"
)

const youdaoApiUrl = "https://openapi.youdao.com/api"

type YouDaoTranslate struct {
	AppId  string
	Secret string
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
	return data.Get("translation").GetIndex(0).String()
}

func TranslateYoudao(from, to, q string) (string, error) {
	return youDaoTranslateStd.Translate(from, to, q)
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

var youDaoTranslateStd = &YouDaoTranslate{
	AppId:  "55ef1f66b723fcf5",
	Secret: "38D3HGvTeMqlQ9M53KVh4YmcsZp2hgk0",
}
