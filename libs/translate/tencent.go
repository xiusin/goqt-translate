package translate

import (
	"crypto/md5"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tmt "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tmt/v20180321"
)

const tencentApiUrl = "tmt.tencentcloudapi.com"
const apiRegion = "ap-beijing"

type TencentTranslate struct {
	AppId  string
	Secret string
}

func TranslateTencent(from, to, q string) (string, error) {
	return tencentTranslateStd.Translate(from, to, q)
}

func (b *TencentTranslate) Translate(from, to, q string) (string, error) {
	credential := common.NewCredential(tencentTranslateStd.AppId, tencentTranslateStd.Secret)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = tencentApiUrl
	client, _ := tmt.NewClient(credential, apiRegion, cpf)

	textRequest := tmt.NewTextTranslateRequest()

	textRequest.SourceText = common.StringPtr(q)
	textRequest.Source = common.StringPtr(from)
	textRequest.Target = common.StringPtr(to)
	textRequest.ProjectId = common.Int64Ptr(0)

	response, err := client.TextTranslate(textRequest)
	if err != nil {
		return "", err
	}
	data, err := simplejson.NewJson([]byte(response.ToJsonString()))
	if err != nil {
		return "", err
	}
	return data.Get("Response").Get("TargetText").String()
}

func (b *TencentTranslate) Sign(salt, q string) string {
	has := md5.Sum([]byte(b.AppId + q + salt + b.Secret))
	return fmt.Sprintf("%x", has)
}

var tencentTranslateStd = &TencentTranslate{
	AppId:  "AKIDRmCdx6p13vfFdtiqlW6wd8RYUGuMBliT",
	Secret: "jHk4Dhz74pRuEPYWI4iibQ4Yu5966wb5",
}
