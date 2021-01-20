package translate

import (
	"testing"
)

func TestTencentTranslate(t *testing.T) {
	t.Log(tencentTranslateStd.Translate("zh", "en", `联通权益卡
`))
}
