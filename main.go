package main

import (
	"context"
	"goqt-translate/libs/components"
	"goqt-translate/libs/helper"
	"goqt-translate/translate"
	"math/rand"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"github.com/xiusin/logger"
)

func init() {
	loggerName := filepath.Join(helper.AppDirPath("logger.log"))
	f, _ := os.OpenFile(loggerName, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	logger.SetOutput(f)
	logger.SetLogLevel(logger.DebugLevel)
	logger.SetReportCaller(true)
	rand.Seed(time.Now().UnixNano())
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("主函数异常：", string(debug.Stack()))
			os.Exit(0)
		}
	}()
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)
	app.SetWindowIcon(gui.NewQIcon5("qrc:/qml/qrc/youdao.png"))
	app.SetStyle(widgets.QStyleFactory_Create("Funsion"))
	app.SetQuitOnLastWindowClosed(false)

	flag := core.Qt__Tool | core.Qt__FramelessWindowHint | core.Qt__X11BypassWindowManagerHint
	//这样新建的窗口在taskbar没有对应的任务图标，并且不 nTopHint | Qt::X11BypassWindowManagerHint);
	window := widgets.NewQMainWindow(nil, flag) // 无边框

	trans := translate.NewTranslateUI(app)
	components.InitSysTray(context.Background(),[]components.MenuAction{
		{
			Text: "轻翻译",
			Callback: func(checked bool) {
				trans.Show()
			},
		},
		{
			Text: "画图板",
			Callback: func(checked bool) {
				components.InitDraw(app)
			},
		},
	}, window)
	app.Exec()
}
