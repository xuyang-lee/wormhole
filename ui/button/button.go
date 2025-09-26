package button

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gorilla/websocket"
	"github.com/xuyang-lee/wormhole/cmd/client"
	"github.com/xuyang-lee/wormhole/hole"
	"github.com/xuyang-lee/wormhole/ui/common"
	"time"
)

func SendText(app fyne.App, input *widget.Entry, messageList *fyne.Container, msgVScroll *container.Scroll) common.Func {
	return func() {
		msg := input.Text
		if msg == "" {
			return
		}

		link, err := hole.GetRandLink()
		if err != nil {
			common.AddSystemMessage(messageList, "can not get connection")
			msgVScroll.ScrollToBottom()
			return
		}

		err = link.Send(websocket.TextMessage, []byte(msg))
		if err != nil {
			common.AddSystemMessage(messageList, fmt.Sprintf("send msg err: %v", err.Error()))
			msgVScroll.ScrollToBottom()
		}

		common.AddSendMessageWithCopyButton(app, messageList, msg)
		msgVScroll.ScrollToBottom()

		input.SetText("") // 清空输入框
	}
}

func ClearMessages(messageList *fyne.Container) common.Func {
	return func() {
		messageList.RemoveAll()
	}
}

func LinkWindow(app fyne.App) common.Func {
	return func() {
		if client.Linked { //如果链接成功,再次点击
			w := app.NewWindow("Link! input the address!")
			out := widget.NewLabel("your wormhole has linked one")
			w.SetContent(out)
			w.Resize(fyne.NewSize(400, 100))
			w.Show()
			time.AfterFunc(time.Second, func() {
				fyne.Do(func() { w.Close() })
			})
			return
		}

		w := app.NewWindow("Link! input the address!")

		input := widget.NewEntry()
		input.SetPlaceHolder("input address here!")

		linkBtn := widget.NewButton("Link!", Link(w, input))

		w.SetContent(container.NewVBox(input, linkBtn))
		w.Resize(fyne.NewSize(400, 100))
		w.Show()
	}
}

func Link(w fyne.Window, input *widget.Entry) common.Func {
	return func() {
		addr := input.Text

		output := widget.NewEntry()
		if err := client.Dial(addr); err != nil {
			output.SetText("link failed... try again... err: " + err.Error())
		} else {
			output.SetText("success! 1s auto close")
		}
		fyne.Do(func() {
			w.SetContent(output)
		})

		time.AfterFunc(time.Second, func() {
			fyne.Do(func() { w.Close() })
		})

	}
}

func CopyAddr(app fyne.App, msg string) common.Func {
	return common.SetClipboard(app, msg)
}
