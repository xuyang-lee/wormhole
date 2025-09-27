package common

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"
)

type Func func()

// direction
const (
	DirectionSystem = iota + 1
	DirectionSend
	DirectionReceive
)

var CurLinkKey string

func GetDirection(d int) string {
	switch d {
	case DirectionSystem:
		return "system"
	case DirectionSend:
		return "send"
	case DirectionReceive:
		return "receive"
	}
	return "unknown"
}

func AddMessage(messageList *fyne.Container, msg string, direction int) {
	label := widget.NewLabel(msg)
	label.Wrapping = fyne.TextWrapWord

	titleLine := fmt.Sprintf("%s: %s", time.Now().Format(time.TimeOnly), GetDirection(direction))
	t := canvas.NewText(titleLine, color.RGBA{0, 255, 255, 255})
	t.TextSize = 12

	// 每条消息一行：Label + Button
	row := container.NewBorder(
		t, nil, nil, nil,
		label,
	)
	messageList.Add(row)
}

func AddSystemMessage(messageList *fyne.Container, msg string) {
	AddMessage(messageList, msg, DirectionSystem)
}

func AddMessageWithCopyButton(app fyne.App, messageList *fyne.Container, msg string, direction int) {
	label := widget.NewLabel(msg)
	label.Wrapping = fyne.TextWrapWord

	titleLine := fmt.Sprintf("%s: %s", time.Now().Format(time.TimeOnly), GetDirection(direction))
	t := canvas.NewText(titleLine, color.RGBA{0, 255, 255, 255})
	t.TextSize = 12

	// 复制按钮，只复制当前消息
	copyBtn := widget.NewButton("copy row", SetClipboard(app, msg))

	// 每条消息一行：Label + Button
	row := container.NewBorder(
		t, nil, nil, copyBtn,
		label,
	)
	messageList.Add(row)
}

func AddSendMessageWithCopyButton(app fyne.App, messageList *fyne.Container, msg string) {
	AddMessageWithCopyButton(app, messageList, msg, DirectionSend)
}

func AddReceiveMessageWithCopyButton(app fyne.App, messageList *fyne.Container, msg string) {
	AddMessageWithCopyButton(app, messageList, msg, DirectionReceive)
}

func SetClipboard(app fyne.App, msg string) Func {
	return func() {
		app.Clipboard().SetContent(msg)
	}
}

func ScrollToBottom(msgVScroll *container.Scroll) {
	fyne.Do(func() {
		msgVScroll.ScrollToBottom()
	})
}
