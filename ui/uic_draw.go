package ui

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/webengine"
	"github.com/therecipe/qt/widgets"
)

type __drawmainwindow struct{}

func (*__drawmainwindow) init() {}

type DrawMainWindow struct {
	*__drawmainwindow
	*widgets.QMainWindow
	Centralwidget *widgets.QWidget
	DrawBoard     *webengine.QWebEngineView
}

func NewDrawMainWindow(p widgets.QWidget_ITF) *DrawMainWindow {
	var par *widgets.QWidget
	if p != nil {
		par = p.QWidget_PTR()
	}
	w := &DrawMainWindow{QMainWindow: widgets.NewQMainWindow(par, 0)}
	w.setupUI()
	w.init()
	return w
}
func (w *DrawMainWindow) setupUI() {
	if w.ObjectName() == "" {
		w.SetObjectName("DrawMainWindow")
	}
	w.Resize2(800, 600)
	if true {
		w.SetStatusTip("")
	}
	w.Centralwidget = widgets.NewQWidget(w, 0)
	w.Centralwidget.SetObjectName("centralwidget")
	w.DrawBoard = webengine.NewQWebEngineView(w.Centralwidget)
	w.DrawBoard.SetObjectName("DrawBoard")
	w.DrawBoard.SetGeometry(core.NewQRect4(0, 0, 801, 601))
	w.DrawBoard.SetAcceptDrops(false)
	w.DrawBoard.SetUrl(core.NewQUrl3("http://localhost:3000/", 0))
	w.SetCentralWidget(w.Centralwidget)
	w.retranslateUi()
	core.QMetaObject_ConnectSlotsByName(w)

}
func (w *DrawMainWindow) retranslateUi() {
	w.SetWindowTitle(core.QCoreApplication_Translate("DrawMainWindow", "MainWindow", "", 0))

}
