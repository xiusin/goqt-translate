package translate

import (
	"crypto/md5"
	"errors"
	"fmt"
	"goqt-translate/libs/request"

	request2 "github.com/mozillazg/request"
	"time"
)

const baiduApiUrl = "http://api.fanyi.baidu.com/api/trans/vip/translate"

type BaiDuTranslate struct {
	AppId  string
	Secret string
}

func TranslateBaidu(from, to, q string) (string, error) {
	return baiduTranslateStd.Translate(from, to, q)
}

func (b *BaiDuTranslate) Translate(from, to, q string) (string, error) {
	salt := fmt.Sprint(time.Now().Unix())
	sign := b.Sign(salt, q)
	req := request.RequestPool.Get()
	defer request.RequestPool.Put(req)
	r := req.(*request2.Request)
	r.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	r.Data = map[string]string{"appid": b.AppId, "q": q, "from": from, "to": to, "salt": salt, "sign": sign}
	resp, err := r.Post(baiduApiUrl) // , u.Encode()
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
		t, _ := resp.Text()
		return "", errors.New(msg + ":" + t + ": 数据: " + fmt.Sprintf("%+v", r.Data))
	}
	m := data.Get("trans_result").GetIndex(0).Get("dst")
	return m.String()
}

func (b *BaiDuTranslate) Sign(salt, q string) string {
	has := md5.Sum([]byte(b.AppId + q + salt + b.Secret))
	return fmt.Sprintf("%x", has)
}

var baiduTranslateStd = &BaiDuTranslate{
	AppId:  "20190318000278289",
	Secret: "AhqLMcMmoEJZaikmc6nM",
}
