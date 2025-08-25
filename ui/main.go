package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/xuyang-lee/wormhole/ui/button"
	"image/color"
)

var address string

func main() {
	// 创建应用
	a := app.New()
	w := a.NewWindow("P2P Sender")

	// title 信息及复制按钮
	address = "ssss:port" //写一个假的，websocket实现了替换成真的
	title := canvas.NewText(address, color.RGBA{0, 255, 0, 255})
	title.TextSize = 12
	addrCopyBtn := widget.NewButton("copy addr", button.CopyAddr(a, address))
	linkBtn := widget.NewButton("link", button.LinkWindow(a))
	l := widget.NewLabel("被动链接:addr发给别人;主动链接:点link,需要别人的addr")

	// 输入框
	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("输入要发送的消息...")

	// 显示日志
	messageList := container.NewVBox()
	msgVScroll := container.NewVScroll(messageList)

	// 按钮
	sendTxtBtn := widget.NewButton("send", button.SendText(a, input, messageList, msgVScroll))
	sendFileBtn := widget.NewButton("send file", button.SendText(a, input, messageList, msgVScroll))
	clearBtn := widget.NewButton("clear", button.ClearMessages(messageList))

	sendBtn := container.NewGridWithColumns(3, sendFileBtn, sendTxtBtn, clearBtn)
	// 界面布局（竖直排列）
	content := container.NewBorder(
		container.NewBorder(title, nil, nil, container.NewGridWithColumns(2, linkBtn, addrCopyBtn), l),
		container.NewVBox(sendBtn, input),
		nil, nil,
		msgVScroll,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 500)) // 初始窗口大小
	w.ShowAndRun()

}

func addMessage(app fyne.App, messages *fyne.Container, msg string) {
	label := widget.NewLabel(msg)
	label.Wrapping = fyne.TextWrapWord

	// 复制按钮，只复制当前消息
	copyBtn := widget.NewButton("复制", func() {
		app.Clipboard().SetContent(msg)
	})

	// 每条消息一行：Label + Button
	row := container.NewHBox(label, copyBtn)
	messages.Add(row)
}
