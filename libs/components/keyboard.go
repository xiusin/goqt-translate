package components

import (
	"fmt"
	hook "github.com/robotn/gohook"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"github.com/xiusin/logger"
	"goqt-translate/libs/request"
	"sort"
	"strings"
)

// 桌面显示当前按键信息
func keyEventHandle(s ...string) {
	hook.Register(hook.KeyDown, s, func(e hook.Event) {
		sort.Strings(s)
		fmt.Println(s)
		globalLabel.SetText(strings.ToUpper(strings.Join(s, " + ")))
		globalLabel.AdjustSize()
	})
}

var globalLabel *widgets.QLabel

func InitKeyboard(app *widgets.QApplication, win *widgets.QMainWindow) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	mm := app.PrimaryScreen().AvailableGeometry()

	win.SetVisible(true)
	win.SetGeometry2(175, mm.Height()-300, 175, 75)
	win.SetAttribute(core.Qt__WA_TranslucentBackground, true)
	win.SetAutoFillBackground(false)
	globalLabel = widgets.NewQLabel(win, core.Qt__FramelessWindowHint|core.Qt__Tool)
	globalLabel.SetFixedHeight(80)

	globalLabel.SetStyleSheet("QLabel{ color: #fff; border: 2px solid #fff; background: rgba(0,0,0,0.5); font-size: 75px; line-height: 80px; }")
	globalLabel.Show()

	modifiers := []string{"ctrl", "shift", "alt", "cmd"}
	for s := range hook.Keycode {
		if request.InArrayStr(s, modifiers) || s == "control" {
			continue
		}
		keyEventHandle(s)
		for i := 0; i < len(modifiers); i++ {
			keyEventHandle(s, modifiers[i])
			for j := i + 1; j < len(modifiers); j++ {
				keyEventHandle(s, modifiers[i], modifiers[j])
				for k := j + 1; k < len(modifiers); k++ {
					keyEventHandle(s, modifiers[i], modifiers[j], modifiers[k])
				}
			}
		}
	}
	win.Show()
	go func() {
		s := hook.Start()
		<-hook.Process(s)
	}()
}
