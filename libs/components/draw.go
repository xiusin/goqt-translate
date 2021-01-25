package components

import (
	"goqt-translate/ui"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"sync"
)
var drawWin *ui.DrawMainWindow
var drawOnce sync.Once

func InitDraw(app *widgets.QApplication) {
	drawOnce.Do(func() {
		drawWin =  ui.NewDrawMainWindow(nil)
		flag := core.Qt__FramelessWindowHint | core.Qt__X11BypassWindowManagerHint
		drawWin.SetWindowFlags(flag)
		//窗口尺寸
		mm := app.PrimaryScreen().AvailableGeometry()
		drawWin.SetGeometry2(0, 0, mm.Width(), mm.Height())
		drawWin.DrawBoard.SetFixedSize2(mm.Width(), mm.Height())
		drawWin.SetAttribute(core.Qt__WA_TranslucentBackground, true)
		drawWin.DrawBoard.Page().SetBackgroundColor(gui.NewQColor2(core.Qt__transparent))
		drawWin.DrawBoard.Show()
	})
	drawWin.Show()
}
