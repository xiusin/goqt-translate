package components

import (
	"fmt"
	"strings"
	"sync"
	"time"

	hook "github.com/robotn/gohook"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"github.com/xiusin/logger"
)

func init() {
	// 定时策略， 150ms内如果有再次按下则直接拼接内容
	// ⇧ ⌃ ⌥ ⌘

}

// 桌面显示当前按键信息
func keyEventHandle(s ...string) {
	hook.Register(hook.KeyDown, s, func(e hook.Event) {
		fmt.Println(s)
		// 每次按下重置timer
		keyboardListenerInstance.Lock()
		defer keyboardListenerInstance.Unlock()

		// 处理10ms的间隔
		if time.Now().Sub(keyboardListenerInstance.prevEnterTime) < 10*time.Millisecond {
			return
		}
		keyboardListenerInstance.prevEnterTime = time.Now()
		if len(s[0]) > 1 {
			return
		}
		keyboardListenerInstance.timer.Reset(2000 * time.Millisecond)
		keyboardListenerInstance.keyStringBuf.WriteString(s[0])
		keyboardListenerInstance.globalLabel.SetText(keyboardListenerInstance.keyStringBuf.String())
	})
}

type keyBoardListener struct {
	globalLabel   *widgets.QLabel
	timer         *time.Timer
	prevEnterTime time.Time
	sync.Mutex
	sync.Once
	win          *widgets.QMainWindow
	keyStringBuf strings.Builder
}

var keyboardListenerInstance = &keyBoardListener{}

func InitKeyboard(app *widgets.QApplication, win *widgets.QMainWindow) {
	keyboardListenerInstance.Do(func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(err)
			}
		}()
		keyboardListenerInstance.timer = time.AfterFunc(2000*time.Millisecond, func() {
			keyboardListenerInstance.Lock()
			defer keyboardListenerInstance.Unlock()
			keyboardListenerInstance.keyStringBuf.Reset() // 清空内容
			// 创建新的label
		})
		keyboardListenerInstance.win = win
		mm := app.PrimaryScreen().AvailableGeometry()
		win.SetVisible(true)
		win.SetGeometry2(120, mm.Height()-300, 575, 30)
		win.SetAttribute(core.Qt__WA_TranslucentBackground, true)
		win.SetAutoFillBackground(false)
		keyboardListenerInstance.globalLabel = widgets.NewQLabel(win, core.Qt__FramelessWindowHint|core.Qt__Tool)
		keyboardListenerInstance.globalLabel.SetFixedHeight(30)
		keyboardListenerInstance.globalLabel.SetFixedWidth(400)
		keyboardListenerInstance.globalLabel.SetStyleSheet("QLabel{ color: #fff; border: 2px solid #fff; background: rgba(0,0,0,0.5); font-size: 30px; }")
		keyboardListenerInstance.globalLabel.Show()

		for s := range hook.Keycode {
			keyEventHandle(s)
		}
		win.Show()
		go func() {
			s := hook.Start()
			<-hook.Process(s)
		}()
	})
}
