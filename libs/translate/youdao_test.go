package translate

import (
	"testing"
)

func TestYoudaoTranslate(t *testing.T) {
	t.Log(youDaoTranslateStd.Translate("en", "zh-CHS", "hello world"))
}
