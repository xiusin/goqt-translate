package components

import (
	"context"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"github.com/xiusin/logger"
	"goqt-translate/libs/helper"
	"runtime/debug"
)

var tray *widgets.QSystemTrayIcon

func InitSysTray(ctx context.Context, menuActions []MenuAction, win *widgets.QMainWindow) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err, string(debug.Stack()))
		}
	}()
	tray = widgets.NewQSystemTrayIcon(win)
	tray.SetIcon(gui.NewQIcon5(":/qml/qrc/logo.png"))

	menus := widgets.NewQMenu(win)
	menuSs, err := helper.GetFileContent("qss/menu.css")
	if err == nil {
		menus.SetStyleSheet(string(menuSs))
	}
	initMenus(win, menuActions, menus)
	tray.SetContextMenu(menus)
	tray.ShowMessage("托盘标题", "托盘显示内容", widgets.QSystemTrayIcon__Information, 1000)
	tray.Show()
}

func GetTray() *widgets.QSystemTrayIcon {
	return tray
}

type MenuAction struct {
	Text     string // 菜单名称
	Icon     string
	ToolTips string
	Callback func(checked bool)
}

type MenuConfig struct {
	Text   string `json:"text"`
	Url    string `json:"url"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
	Flag   string `json:"flag"`
}

var menuActionArr []*widgets.QAction

func initMenus(win *widgets.QMainWindow, menuActions []MenuAction, menus *widgets.QMenu) *widgets.QMenu {
	for _, action := range menuActions {
		menuAction := menus.AddAction(action.Text)
		menuAction.SetParent(win)
		menuAction.SetToolTip(action.ToolTips)
		menuAction.ConnectTriggered(action.Callback)
		if action.Icon != "" {
			menuAction.SetIcon(gui.NewQIcon5(action.Icon))
		}
		menuActionArr = append(menuActionArr, menuAction)
	}
	return menus
}
