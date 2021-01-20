package notification

import (
	gosxnotifier "github.com/deckarep/gosx-notifier"
)

func Notification(title, content string) error {
	//At a minimum specifiy a message to display to end-user.
	note := gosxnotifier.NewNotification(content)

	//Optionally, set a title
	note.Title = title

	//Optionally, set a subtitle
	//note.Subtitle = "每日花钱时刻"

	//Optionally, set a sound from a predefined set.
	note.Sound = gosxnotifier.Basso

	//Optionally, set a group which ensures only one notification is ever shown replacing previous notification of same group id.
	//note.Group = "com.unique.yourapp.identifier"

	//Optionally, set a sender (Notification will now use the Safari icon)
	//note.Sender = "com.apple.Safari"

	//Optionally, specifiy a url or bundleid to open should the notification be
	//clicked.
	//note.Link = "http://www.yahoo.com" //or BundleID like: com.apple.Terminal

	//Optionally, an app icon (10.9+ ONLY)
	//note.AppIcon = "fund.png"

	//Optionally, a content image (10.9+ ONLY)
	note.ContentImage = "fund.png"

	//Then, push the notification
	return note.Push()
}
