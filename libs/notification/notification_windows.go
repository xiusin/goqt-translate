package notification

import (
	"github.com/go-toast/toast"
)

func Notification(title, content string) error {
	notification := toast.Notification{
		AppID:   "Microsoft.Windows.Shell.RunDialog",
		Title:   title,
		Message: content,
		Icon:    "C:\\path\\to\\your\\logo.png", // 文件必须存在
		//Actions: []toast.Action{
		//	{"protocol", "按钮1", "https://www.google.com/"},
		//	{"protocol", "按钮2", "https://github.com/"},
		//},
	}
	return notification.Push()
}
