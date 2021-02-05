package components

import (
	"fmt"
	"goqt-translate/libs/sound"
	"strings"
	"sync"
	"time"

	hook "github.com/robotn/gohook"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"github.com/xiusin/logger"
)

var keyMap = map[string]string{
	"shift":     "⇧",
	"lshift":    "⇧",
	"rshift":    "⇧",
	"command":   "⌘",
	"rcmd":      "⌘",
	"lcmd":      "⌘",
	"cmd":       "⌘",
	"control":   "⌃",
	"rctrl":     "⌃",
	"lctrl":     "⌃",
	"ctrl":      "⌃",
	"enter":     "↩",
	"option":    "⌥",
	"caps lock": "⇪",
	"return":    "↩",
	"esc":       "⎋",
	"delete":    "⌫",
	"up":        "↑",
	"down":      "↓",
	"right":     "→",
	"left":      "←",
	"page up":   "⇞",
	"page down": "⇟",
	"tab":       "⇥",
}

var raw2key = map[uint16]string{ // https://github.com/wesbos/keycodes
	0:   "error",
	3:   "break",
	8:   "backspace",
	9:   "tab",
	12:  "clear",
	13:  "enter",
	16:  "shift",
	17:  "ctrl",
	18:  "alt",
	19:  "pause/break",
	20:  "caps lock",
	21:  "hangul",
	25:  "hanja",
	27:  "escape",
	28:  "conversion",
	29:  "non-conversion",
	32:  "spacebar",
	33:  "page up",
	34:  "page down",
	35:  "end",
	36:  "home",
	37:  "left arrow",
	38:  "up arrow",
	39:  "right arrow",
	40:  "down arrow",
	41:  "select",
	42:  "print",
	43:  "execute",
	44:  "Print Screen",
	45:  "insert",
	46:  "delete",
	47:  "help",
	48:  "0",
	49:  "1",
	50:  "2",
	51:  "3",
	52:  "4",
	53:  "5",
	54:  "6",
	55:  "7",
	56:  "8",
	57:  "9",
	58:  ":",
	59:  ";",
	60:  "<",
	61:  "=",
	63:  "ß",
	64:  "@",
	65:  "a",
	66:  "b",
	67:  "c",
	68:  "d",
	69:  "e",
	70:  "f",
	71:  "g",
	72:  "h",
	73:  "i",
	74:  "j",
	75:  "k",
	76:  "l",
	77:  "m",
	78:  "n",
	79:  "o",
	80:  "p",
	81:  "q",
	82:  "r",
	83:  "s",
	84:  "t",
	85:  "u",
	86:  "v",
	87:  "w",
	88:  "x",
	89:  "y",
	90:  "z",
	91:  "l-super",
	92:  "r-super",
	93:  "apps",
	95:  "sleep",
	96:  "numpad 0",
	97:  "numpad 1",
	98:  "numpad 2",
	99:  "numpad 3",
	100: "numpad 4",
	101: "numpad 5",
	102: "numpad 6",
	103: "numpad 7",
	104: "numpad 8",
	105: "numpad 9",
	106: "multiply",
	107: "add",
	108: "numpad period",
	109: "subtract",
	110: "decimal point",
	111: "divide",
	112: "f1",
	113: "f2",
	114: "f3",
	115: "f4",
	116: "f5",
	117: "f6",
	118: "f7",
	119: "f8",
	120: "f9",
	121: "f10",
	122: "f11",
	123: "f12",
	124: "f13",
	125: "f14",
	126: "f15",
	127: "f16",
	128: "f17",
	129: "f18",
	130: "f19",
	131: "f20",
	132: "f21",
	133: "f22",
	134: "f23",
	135: "f24",
	144: "num lock",
	145: "scroll lock",
	160: "^",
	161: "!",
	162: "؛",
	163: "#",
	164: "$",
	165: "ù",
	166: "page backward",
	167: "page forward",
	168: "refresh",
	169: "closing paren (AZERTY)",
	170: "*",
	171: "~ + * key",
	172: "home key",
	173: "minus (firefox), mute/unmute",
	174: "decrease volume level",
	175: "increase volume level",
	176: "next",
	177: "previous",
	178: "stop",
	179: "play/pause",
	180: "e-mail",
	181: "mute/unmute (firefox)",
	182: "decrease volume level (firefox)",
	183: "increase volume level (firefox)",
	186: "semi-colon / ñ",
	187: "equal sign",
	188: "comma",
	189: "dash",
	190: "period",
	191: "forward slash / ç",
	192: "grave accent / ñ / æ / ö",
	193: "?, / or °",
	194: "numpad period (chrome)",
	219: "open bracket",
	220: "back slash",
	221: "close bracket / å",
	222: "single quote / ø / ä",
	223: "`",
	224: "left or right ⌘ key (firefox)",
	225: "altgr",
	226: "< /git >, left back slash",
	230: "GNOME Compose Key",
	231: "ç",
	233: "XF86Forward",
	234: "XF86Back",
	235: "non-conversion",
	240: "alphanumeric",
	242: "hiragana/katakana",
	243: "half-width/full-width",
	244: "kanji",
	251: "unlock trackpad (Chrome/Edge)",
	255: "toggle touchpad",
}

type keyBoardListener struct {
	globalLabel     *widgets.QLabel
	timer           *time.Timer
	prevEnterTime   map[string]time.Time
	ModifierMapping map[string]bool
	sync.Mutex
	sync.Once
	keyCodeMapKeyChar map[uint16]string
	win               *widgets.QMainWindow
	keyStringBuf      strings.Builder
}

var keyboardListenerInstance = &keyBoardListener{}

func InitKeyboard(app *widgets.QApplication) *keyBoardListener { //, win *widgets.QMainWindow
	keyboardListenerInstance.Do(func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(err)
			}
		}()
		keyboardListenerInstance.prevEnterTime = map[string]time.Time{}
		keyboardListenerInstance.keyCodeMapKeyChar = map[uint16]string{}
		keyboardListenerInstance.ModifierMapping = map[string]bool{}

		keyboardListenerInstance.timer = time.AfterFunc(2000*time.Millisecond, func() {
			keyboardListenerInstance.Lock()
			defer keyboardListenerInstance.Unlock()
			keyboardListenerInstance.keyStringBuf.Reset() // 清空内容
			// 创建新的label
		})
		win := widgets.NewQWidget(nil, core.Qt__FramelessWindowHint|core.Qt__WindowStaysOnTopHint)
		win.SetVisible(true)
		win.SetAttribute(core.Qt__WA_TranslucentBackground, true)
		win.SetAutoFillBackground(false)
		//painter := gui.NewQPainter()
		//win.InitPainter(painter)
		//win.ConnectPaintEvent(func(event *gui.QPaintEvent) {
		//	painter.SetPen3(core.Qt__NoPen)
		//	painter.SetBrush(gui.NewQBrush4(core.Qt__red, 0))
		//	painter.SetRenderHint(gui.QPainter__Antialiasing, true)
		//	rect := win.Rect()
		//	rect.SetWidth(rect.Width() - 1)
		//	rect.SetHeight(rect.Height() - 1)
		//	painter.DrawRoundedRect3(rect, 10, 10, 0)
		//})

		keyboardListenerInstance.globalLabel = widgets.NewQLabel(nil, core.Qt__FramelessWindowHint|core.Qt__WindowStaysOnTopHint)
		keyboardListenerInstance.globalLabel.SetScaledContents(true)
		keyboardListenerInstance.globalLabel.SetStyleSheet("QLabel{ color: #fff; border: 2px solid #fff; border-radius: 0wwwwwwwwwwwwpx; background: rgba(0,0,0,0.4); font-size: 30px; }")
		keyboardListenerInstance.globalLabel.SetGeometry2(15, app.PrimaryScreen().AvailableSize().Height()-35, 0, 0)

		for s, v := range hook.Keycode {
			keyboardListenerInstance.keyCodeMapKeyChar[v] = s
		}
		player, soundData := sound.InitStreamer()
		go func() {
			s := hook.Start()
			for ev := range s {
				var kc string
				if ev.Keychar != 65535 {
					kc = string(ev.Keychar)
				} else {
					kc = keyboardListenerInstance.keyCodeMapKeyChar[ev.Keycode]
				}
				if ev.Kind == hook.KeyUp {
					delete(keyboardListenerInstance.ModifierMapping, kc)
				} else if ev.Kind == hook.KeyHold { // ev.Kind == hook.KeyHold || 修饰键用keyHold
					fmt.Println(ev)
					func(keyChar string) {
						keyboardListenerInstance.Lock()
						defer keyboardListenerInstance.Unlock()
						go func() {
							p := player.NewPlayer()
							p.Write(soundData)
							p.Close()
						}()
						keyboardListenerInstance.prevEnterTime[keyChar] = time.Now()
						if len(keyChar) > 1 {
							keyboardListenerInstance.ModifierMapping[kc] = true
							if t, ok := keyMap[keyChar]; ok {
								keyChar = t
							} else {
								if keyChar == "space" {
									keyChar = " "
								}
							}
						}
						keyboardListenerInstance.timer.Reset(2000 * time.Millisecond)
						keyboardListenerInstance.keyStringBuf.WriteString(keyChar)
					}(kc)
				}
			}
			<-hook.Process(s)
		}()

		qtimer := core.NewQTimer(nil)

		// 记录上次的宽高,如果未发生变化则不执行
		var prevWidth, prevHeight int

		qtimer.ConnectTimeout(func() {
			if keyboardListenerInstance.keyStringBuf.Len() == 0 {
				prevHeight, prevWidth = 0, 0
				win.Hide()
				keyboardListenerInstance.globalLabel.Hide()
			} else {
				win.Show()
				keyboardListenerInstance.globalLabel.Show()
			}
			keyboardListenerInstance.globalLabel.SetText(keyboardListenerInstance.keyStringBuf.String())

			keyboardListenerInstance.globalLabel.AdjustSize()
			w := keyboardListenerInstance.globalLabel.Width()
			h := keyboardListenerInstance.globalLabel.Height()
			if w == prevWidth && h == prevHeight {
				return
			}
			prevWidth, prevHeight = w, h
			win.SetFixedSize2(w, h)
		})
		qtimer.SetInterval(50)
		qtimer.Start2()
	})
	return keyboardListenerInstance
}
