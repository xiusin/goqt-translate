package helper

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"unicode"
)

var appleScriptFile string

func init() {
	if runtime.GOOS == "darwin" { // darwin执行自己的文件
		tmp, _ := ioutil.TempFile("", "")
		appleScriptFile = tmp.Name()
		tmp.WriteString(osaScript())
		tmp.Sync() // 写内容
	}
}

// GetFileContent 获取文件内容
func GetFileContent(name string) ([]byte, error) {
	return ioutil.ReadFile(AppDirPath(name))
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

// getCurrentDir 获取当前文件夹
func getCurrentDir() string {
	workingDir, _ := os.Getwd()
	if runtime.GOOS == "darwin" {
		workingDir, _ = os.Executable()
		workingDir = filepath.Dir(workingDir)
	}
	return workingDir
}

func AppDirPath(name string) string {
	return filepath.Join(getCurrentDir(), name)
}

// ExecToCopy 执行apple script脚本
func ExecToCopy() error{
	cmd := exec.Command("osascript", appleScriptFile)
	return cmd.Run()
}

// osaScript 临时复制方案
func osaScript() string {
	return `tell application "System Events"
  key code 8 using command down
end tell`
}
