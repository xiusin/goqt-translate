package translate

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	request2 "github.com/mozillazg/request"
	"goqt-translate/libs/request"
	"time"
)

const XfTransDomain = "itrans.xfyun.cn"

const xfTransUri = "/v2/its"

const xfTransUrl = "https://" + XfTransDomain + xfTransUri

type XunFeiTranslate struct {
	AppId  string
	ApiKey string
	Secret string
}

func (b *XunFeiTranslate) Translate(from, to, q string) (string, error) {
	req := request.RequestPool.Get()
	defer request.RequestPool.Put(req)
	r := req.(*request2.Request)
	data := map[string]interface{}{
		"common": map[string]string{
			"app_id": b.AppId,
		},
		"business": map[string]string{
			"from": from,
			"to":   to,
		},
		"data": map[string]string{
			"text": base64.StdEncoding.EncodeToString([]byte(q)),
		},
	}
	body, err := json.Marshal(&data)
	if err != nil {
		return "", err
	}
	strBody := string(body)
	r.Headers = map[string]string{
		"Content-Type": "application/json",
		"Date":         time.Now().UTC().Format(time.RFC1123),
		"Digest":       "SHA-256=" + signBody(&strBody),
	}
	sign := generateSignature(XfTransDomain, r.Headers["Date"], "POST", xfTransUri, "HTTP/1.1", r.Headers["Digest"], b.Secret)
	authHeader := fmt.Sprintf(`api_key="%s", algorithm="%s", headers="host date request-line digest", signature="%s"`,
		b.ApiKey, "hmac-sha256", sign)
	r.Headers["Authorization"] = authHeader

	r.Json = data

	resp, err := r.Post(xfTransUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	resData, err := resp.Json()
	if err != nil {
		return "", err
	}

	code, err := resData.Get("code").Int()
	if err != nil {
		return "", err
	}
	if code != 0 {
		msg, _ := resData.Get("message").String()
		if msg == "" {
			msg, _ = resp.Text()
		}
		return "", errors.New(msg + " 请求数据: " + strBody)
	}
	return resData.Get("data").Get("result").Get("trans_result").Get("dst").String()
}

func TranslateXunfei(from, to, q string) (string, error) {
	return xunFeiTranslateStd.Translate(from, to, q)
}

func (b *XunFeiTranslate) Sign(salt, curTime, q string) string {
	return ""
}

func generateSignature(host, date, httpMethod, requestUri, httpProto, digest string, secret string) string {
	// 不是request-line的话，则以header名称,后跟ASCII冒号:和ASCII空格，再附加header值
	var signatureStr string
	if len(host) != 0 {
		signatureStr = "host: " + host + "\n"
	}
	signatureStr += "date: " + date + "\n"
	// 如果是request-line的话，则以 http_method request_uri http_proto
	signatureStr += httpMethod + " " + requestUri + " " + httpProto + "\n"
	signatureStr += "digest: " + digest
	return hmacsign(signatureStr, secret)
}
func hmacsign(data, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}
func signBody(data *string) string {
	// 进行sha256签名
	sha := sha256.New()
	sha.Write([]byte(*data))
	encodeData := sha.Sum(nil)
	// 经过base64转换
	return base64.StdEncoding.EncodeToString(encodeData)
}

var xunFeiTranslateStd = &XunFeiTranslate{
	AppId:  "5fba6eaa",
	Secret: "053250b6d74c387f98db26aa5814e9b7",
	ApiKey: "4b63a77ae194dedc15555d898209629f",
}
