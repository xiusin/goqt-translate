package components

import (
	"goqt-translate/ui"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"sync"
)

type draw struct {
	*ui.DrawMainWindow
	loaded bool
	drawOnce sync.Once
}
var drawInstance = &draw{}
func InitDraw(app *widgets.QApplication) {
	drawInstance.drawOnce.Do(func() {
		drawInstance.DrawMainWindow = ui.NewDrawMainWindow(nil)
		drawInstance.DrawBoard.Page().SetUrl(core.NewQUrl3("http://127.0.0.1:11731/#/", 0))
		drawInstance.DrawBoard.Page().ConnectLoadFinished(func(ok bool) {
			drawInstance.loaded = true
			drawInstance.Show()
		})
		flag := core.Qt__FramelessWindowHint | core.Qt__X11BypassWindowManagerHint
		drawInstance.SetWindowFlags(flag)
		//窗口尺寸
		mm := app.PrimaryScreen().AvailableGeometry()
		drawInstance.SetGeometry2(0, 0, mm.Width(), mm.Height())
		drawInstance.DrawBoard.SetFixedSize2(mm.Width(), mm.Height())
		drawInstance.SetAttribute(core.Qt__WA_TranslucentBackground, true)
		drawInstance.DrawBoard.Page().SetBackgroundColor(gui.NewQColor2(core.Qt__transparent))
		drawInstance.DrawBoard.Show()
	})
	if drawInstance.IsHidden() && drawInstance.loaded {
		drawInstance.Show()
	} else {
		drawInstance.Hide()
	}
}
