package translate

import (
	"goqt-translate/libs/helper"
	"goqt-translate/libs/translate"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"github.com/xiusin/logger"
)

type TranslateUI struct {
	sync.Mutex
	*widgets.QWidget
	fromInput              *widgets.QTextEdit
	baidu                  *widgets.QTextEdit
	youdao                 *widgets.QTextEdit
	xunfei                 *widgets.QTextEdit
	tencent                *widgets.QTextEdit
	tabs                   *widgets.QTabWidget
	tipsWidget             *widgets.QWidget
	tipsOutWidget          *widgets.QWidget
	btn                    *widgets.QPushButton
	clipboard              *gui.QClipboard
	transing               bool
	requesting             bool
	isDrag                 bool
	isKeyDown              bool
	tipsTickTimer          *time.Ticker
	tipWidgetHideCh        chan struct{}
	shiftClick             uint32
	shiftClickEventCh      chan uint8
	shiftClickCanelEventCh chan struct{}
	soundPlaying           bool

	wg sync.WaitGroup
}

func NewTranslateUI(app *widgets.QApplication) *TranslateUI {
	width := 300
	win := widgets.NewQWidget(nil, core.Qt__WindowStaysOnTopHint)
	win.SetWindowTitle("轻翻译")
	win.SetFixedSize2(width, 600)
	fanyiIcon := gui.NewQIcon5(":/qml/qrc/youdao.png")
	win.SetWindowIcon(fanyiIcon)
	layout := widgets.NewQVBoxLayout()
	layout.SetSpacing(0)
	layout.SetSizeConstraint(widgets.QLayout__SetFixedSize) //设置布局跟随组件的宽高
	layout.AddSpacing(0)
	layout.SetContentsMargins(0, 0, 0, 0)
	win.SetLayout(layout)

	formInput := widgets.NewQTextEdit(win)
	formInput.SetFixedSize2(width, 250)

	baidu := widgets.NewQTextEdit(win)
	baidu.SetFixedSize2(width, 350)
	baidu.SetReadOnly(true)
	baidu.Show()

	youdao := widgets.NewQTextEdit(win)
	youdao.SetFixedSize2(width, 350)
	youdao.SetReadOnly(true)
	youdao.Show()

	xunfei := widgets.NewQTextEdit(win)
	xunfei.SetFixedSize2(width, 350)
	xunfei.SetReadOnly(true)
	xunfei.Show()

	tencent := widgets.NewQTextEdit(win)
	tencent.SetFixedSize2(width, 350)
	tencent.SetReadOnly(true)
	tencent.Show()

	tipsWidget := widgets.NewQWidget(nil, core.Qt__Tool|core.Qt__WindowStaysOnTopHint)
	tipsWidget.SetWindowTitle("翻译结果")
	tipsWidget.SetFixedSize2(width, 351)
	qtabSs, err := helper.GetFileContent("qss/qtab.css")
	// tabs
	tabs := widgets.NewQTabWidget(tipsWidget)
	tabs.SetFixedSize2(width+1, 351)
	tabs.AddTab(youdao, "有道")
	tabs.SetTabIcon(0, fanyiIcon)
	tabs.AddTab(baidu, "百度")
	tabs.SetTabIcon(1, gui.NewQIcon5(":/qml/qrc/baidu.png"))
	tabs.AddTab(tencent, "腾讯")
	tabs.SetTabIcon(2, gui.NewQIcon5(":/qml/qrc/tencent.png"))
	tabs.AddTab(xunfei, "讯飞")
	tabs.SetTabIcon(3, gui.NewQIcon5(":/qml/qrc/xunfei.png"))

	if err == nil {
		tipsWidget.SetStyleSheet(string(qtabSs))
	}

	button1 := widgets.NewQPushButton2("翻译", win)
	button1.SetFixedSize2(width, 30)

	ss, err := helper.GetFileContent("qss/translate.css")
	if err == nil {
		button1.SetStyleSheet(string(ss))
	}
	layout.SetSizeConstraint(widgets.QLayout__SetFixedSize) //设置布局跟随组件的宽高
	layout.AddWidget(formInput, 1, 0)
	layout.AddWidget(button1, 1, 0)
	layout.AddWidget(tipsWidget, 0, 0)

	tu := &TranslateUI{
		QWidget:                win,
		fromInput:              formInput,
		baidu:                  baidu,
		youdao:                 youdao,
		xunfei:                 xunfei,
		tencent:                tencent,
		btn:                    button1,
		tabs:                   tabs,
		tipsWidget:             tipsWidget,
		tipsOutWidget:          tipsWidget,
		clipboard:              app.Clipboard(),
		tipsTickTimer:          time.NewTicker(2 * time.Second),
		shiftClick:             0,
		tipWidgetHideCh:        make(chan struct{}),
		shiftClickCanelEventCh: make(chan struct{}),
		shiftClickEventCh:      make(chan uint8),
	}

	soundLabel := widgets.NewQLabel(win, 0) //":/qml/qrc/sound.png"
	pixmap := gui.NewQPixmap3(":/qml/qrc/sound.png", "", core.Qt__AutoColor)

	pixmap = pixmap.Scaled2(32, 32, core.Qt__KeepAspectRatio, core.Qt__SmoothTransformation)

	soundLabel.SetPixmap(pixmap)
	soundLabel.SetGeometry2(240, 330, 35, 35)
	soundLabel.SetToolTip("双击播放声音")
	soundLabel.ConnectMouseDoubleClickEvent(func(event *gui.QMouseEvent) {
		if tu.soundPlaying {
			return
		}
		txt := strings.TrimSpace(tu.fromInput.ToPlainText())
		if len(txt) > 0 {
			toLang := "en"
			if !helper.IsChinese(&txt) {
				toLang = "zh"
			}
			go func() {
				tu.soundPlaying = true
				defer func() {
					tu.soundPlaying = false
				}()
				translate.ToVoice(txt, toLang)
			}()
		}
	})
	soundLabel.Show()

	tu.registerEvent()
	tu.ShowFromButton()
	return tu
}

func (tu *TranslateUI) registerEvent() {
	tu.btn.ConnectClicked(transEvent(tu, tu.btn))
	for _, v := range []*widgets.QTextEdit{tu.baidu, tu.youdao, tu.xunfei} {
		v.ConnectSelectionChanged(func(v *widgets.QTextEdit) func() {
			return func() {
				if to := v.ToPlainText(); to != "" {
					tu.clipboard.SetText(to, gui.QClipboard__Clipboard)
				}
			}
		}(v))
	}

	tu.ConnectCloseEvent(func(event *gui.QCloseEvent) {
		go func() {
			time.Sleep(time.Millisecond * 200)
			tu.ShowFromButton()
		}()
	})

	go tu.listenDoubleShiftClick()
	go tu.listenKeyBoardOrMouseEvent()

}

func (tu *TranslateUI) listenKeyBoardOrMouseEvent() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	s := hook.Start()
	defer hook.End()
	for ev := range s {
		tu.Lock()
		ctrlKeyCode := map[uint16]struct{}{29: {}}
		var altKeyCode uint16 = 56 // alt
		if runtime.GOOS == "darwin" {
			ctrlKeyCode = map[uint16]struct{}{3675: {}, 3676: {}} // option
			altKeyCode = 29
		}

		if _, ok := ctrlKeyCode[ev.Keycode]; ok && ev.Kind == hook.KeyHold {
			if !tu.isKeyDown {
				tu.isKeyDown = true
			}
		}

		// 使用alt键
		if ev.Keycode == altKeyCode && ev.Kind == hook.KeyUp {
			if tu.IsHidden() { // 非展示状态
				if atomic.LoadUint32(&tu.shiftClick) == 0 {
					go func() { tu.shiftClickEventCh <- 1 }()
				} else if atomic.LoadUint32(&tu.shiftClick) == 1 {
					go func() { tu.shiftClickEventCh <- 2 }()
				}
			}
		}

		// todo 处理粘键问题
		if _, ok := ctrlKeyCode[ev.Keycode]; ok && ev.Kind == hook.KeyUp {
			if tu.isKeyDown && tu.isDrag {
				tu.isDrag = false
				old := tu.clipboard.Text(gui.QClipboard__Clipboard) // 记录之前的内容
				if runtime.GOOS == "darwin" {
					if err := helper.ExecToCopy(); err != nil {
						logger.Print("执行命令失败:", err)
						tu.Unlock()
						continue
					}
					time.Sleep(time.Millisecond * 30)
				} else {
					robotgo.KeyTap("c", "ctrl")
				}
				selectionTxt := tu.clipboard.Text(gui.QClipboard__Clipboard)
				tu.clipboard.SetText(old, gui.QClipboard__Clipboard)
				tu.fromInput.SetText(selectionTxt)
				tu.HideFromButton() // 隐藏相关按钮
				tu.btn.Click()
			}
			tu.isKeyDown = false
			tu.isDrag = false
		}
		if ev.Kind == hook.MouseDrag && tu.isKeyDown { // 监听鼠标抬起事件
			if !tu.isDrag {
				tu.isDrag = true
			}
		}
		tu.Unlock()
	}
}

func (tu *TranslateUI) ShowFromButton() {
	tu.fromInput.Show()
	tu.btn.Show()
}

func (tu *TranslateUI) IsFull() bool {
	return tu.fromInput.IsHidden() == false
}

func (tu *TranslateUI) HideFromButton() {
	tu.fromInput.Hide()
	tu.btn.Hide()
}

func (tu *TranslateUI) listenDoubleShiftClick() {
	for {
		select {
		case val := <-tu.shiftClickEventCh:
			if tu.IsFull() {
				if val >= 2 {
					atomic.StoreUint32(&tu.shiftClick, 0)
					tu.ShowFromButton()
					tu.Show()
					tu.fromInput.SetFocus2()
				} else {
					atomic.StoreUint32(&tu.shiftClick, 1)
					go func() {
						time.Sleep(300 * time.Millisecond)
						tu.shiftClickCanelEventCh <- struct{}{}
					}()
				}
			}
		case <-tu.shiftClickCanelEventCh:
			atomic.StoreUint32(&tu.shiftClick, 0)
		}
	}
}

type transType struct {
	Name     string
	From, To string
	Obj      *widgets.QTextEdit
	fun      func(string, string, string) (string, error)
}

func transEvent(tu *TranslateUI, btn *widgets.QPushButton) func(bool) {
	return func(checked bool) {
		if tu.requesting {
			return
		}
		originTxt := btn.Text()
		objs := []*transType{
			{"有道", "en", "zh-CHS", tu.youdao, translate.TranslateYoudao},
			{"百度", "en", "zh", tu.baidu, translate.TranslateBaidu},
			{"腾讯", "en", "zh", tu.tencent, translate.TranslateTencent},
			{"讯飞", "en", "cn", tu.xunfei, translate.TranslateXunfei},
		}
		txt := strings.Trim(tu.fromInput.ToPlainText(), " \n")
		if txt != "" {
			go func() {
				btn.SetText("正在翻译")
				tu.requesting = true
				isChinese := helper.IsChinese(&txt)
				// 转换声音
				tu.wg.Add(len(objs))
				for _, v := range objs {
					go func(v *transType) {
						defer tu.wg.Done()
						if isChinese {
							v.From, v.To = v.To, v.From
						}
						rest, err := v.fun(v.From, v.To, txt)
						if err != nil {
							logger.Error(v.Name, " ", err)
							v.Obj.SetText("翻译异常: " + err.Error())
						} else {
							v.Obj.SetText(rest)
						}
					}(v)
				}
				tu.wg.Wait()
				btn.SetText(originTxt)
				if !tu.IsFull() {
					x, y := robotgo.GetMousePos() // 获取当前鼠标位置并展示弹窗
					tu.SetGeometry2(x, y, 300, 350)
				}
				tu.Show()
				tu.tabs.TabBarClicked(0) // 默认显示第一个
				tu.requesting = false
			}()
		}
	}
}

func (tu *TranslateUI) Show() {
	tu.QWidget.Show()
	if tu.IsMinimized() {
		tu.SetWindowState(core.Qt__WindowActive | tu.WindowState()&^core.Qt__WindowMinimized)
	} else {
		tu.QWidget.SetWindowState(core.Qt__WindowActive)
	}
	tu.SetAttribute(core.Qt__WA_Mapped, true)
	tu.ActivateWindow()
	tu.Raise()
	tu.SetFocus2()
	tu.fromInput.SetFocus2()
}
