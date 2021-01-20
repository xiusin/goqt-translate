package main

import (
	"fmt"
	"goqt-translate/translate"
	"math/rand"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/therecipe/qt/widgets"
	"github.com/xiusin/logger"
)

func init() {
	// _ = os.Remove(helper.UserHomeDir("error.log"))
	// f, _ := os.OpenFile(helper.UserHomeDir("error.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	// logger.SetOutput(f)
	// logger.SetLogLevel(logger.DebugLevel)
	// logger.SetReportCaller(true)
	rand.Seed(time.Now().UnixNano())
}

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	defer func() {
		if err := recover(); err != nil {
			logger.Error("主函数异常：", string(debug.Stack()))
			os.Exit(0)
		}
	}()
	app.SetStyle(widgets.QStyleFactory_Create("Funsion"))
	app.SetQuitOnLastWindowClosed(false)

	trans := translate.NewTranslateUI(app)
	if runtime.GOOS == "darwin" {
		fmt.Println("双击alt打开翻译界面")
		fmt.Println("按下alt选中要翻译的文本放开alt即可弹出翻译内容")
	}
	trans.Show()
	app.Exec()
}
