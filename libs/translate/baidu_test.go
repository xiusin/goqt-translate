package translate

import (
	"testing"
)

func TestBaiDuTranslate(t *testing.T) {
	t.Log(baiduTranslateStd.Translate("en", "zh", `联通权益卡
`))
}
