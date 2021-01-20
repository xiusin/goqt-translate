package helper

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"unicode"
)

// GetFileContent 获取文件内容
func GetFileContent(name string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(getCurrentDir(), name))
}

// IsChinese 判断是否为中文语句
func IsChinese(str *string) bool {
	h, e := 0, 0
	for _, r := range *str {
		if unicode.Is(unicode.Han, r) {
			h++
		} else {
			e++
		}
	}
	return h > e
}

func getCurrentDir() string {
	workingDir, _ := os.Getwd()
	if runtime.GOOS == "darwin" {
		workingDir, _ = os.Executable()
	}
	return workingDir
}
