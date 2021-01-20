package translate

import (
	"testing"
)

func TestXunfeiTranslate(t *testing.T) {
	t.Log(xunFeiTranslateStd.Translate("cn", "en", `联通权益卡
`))
}
