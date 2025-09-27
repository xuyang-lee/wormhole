package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/xuyang-lee/wormhole/config"
	"github.com/xuyang-lee/wormhole/hole"
	"github.com/xuyang-lee/wormhole/ui/button"
	"github.com/xuyang-lee/wormhole/ui/common"
	"github.com/xuyang-lee/wormhole/ui/receive"
	"image/color"
	"time"
)

var address string

func main() {
	// 创建应用
	a := app.New()
	w := a.NewWindow("P2P Sender")

	// 开启后台服务
	go hole.Init()

	time.Sleep(time.Second)
	// title 信息及复制按钮
	address = fmt.Sprintf(":%d", config.Conf.Port)
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
	sendFileBtn := widget.NewButton("send file", button.SendFile(input, messageList, msgVScroll))
	closeBtn := widget.NewButton("close", button.SendClose(messageList, msgVScroll))
	clearBtn := widget.NewButton("clear", button.ClearMessages(messageList))

	btnGrid := container.NewGridWithColumns(4, closeBtn, sendFileBtn, sendTxtBtn, clearBtn)
	// 界面布局（竖直排列）
	content := container.NewBorder(
		container.NewBorder(title, nil, nil, container.NewGridWithColumns(2, linkBtn, addrCopyBtn), l),
		container.NewVBox(btnGrid, input),
		nil, nil,
		msgVScroll,
	)

	listen := common.NewListener(func() {
		common.AddSystemMessage(messageList, "已连接")
		common.ScrollToBottom(msgVScroll)
	})
	// 链接监听，注册监听者
	hole.RegisterListener(listen)
	// 开启接受后台
	go receive.Receive(a, messageList, msgVScroll)
	//go receive.TestReceive(a, messageList, msgVScroll)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 500)) // 初始窗口大小
	w.ShowAndRun()

}
