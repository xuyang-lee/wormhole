package button

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/xuyang-lee/wormhole/cmd/client"
	"image/color"
	"time"
)

type Func func()

func SendText(app fyne.App, input *widget.Entry, messageList *fyne.Container, msgVScroll *container.Scroll) Func {
	return func() {
		msg := input.Text
		if msg == "" {
			return
		}
		// TODO: 调用你自己的 transport.SendText(addr, msg)
		addMessage(app, messageList, msg)
		msgVScroll.ScrollToBottom()

		input.SetText("") // 清空输入框
	}
}

func addMessage(app fyne.App, messageList *fyne.Container, msg string) {
	label := widget.NewLabel(msg)
	label.Wrapping = fyne.TextWrapWord

	t := canvas.NewText(time.Now().Format(time.TimeOnly+": send"), color.RGBA{0, 255, 255, 255})
	t.TextSize = 12

	// 复制按钮，只复制当前消息
	copyBtn := widget.NewButton("copy row", CopyAddr(app, msg))

	// 每条消息一行：Label + Button
	row := container.NewBorder(
		t, nil, nil, copyBtn,
		label,
	)
	messageList.Add(row)
}

func ClearMessages(messageList *fyne.Container) Func {
	return func() {
		messageList.Objects = make([]fyne.CanvasObject, 0, len(messageList.Objects))
	}
}

func CopyAddr(app fyne.App, msg string) Func {
	return func() {
		app.Clipboard().SetContent(msg)
	}
}

func LinkWindow(app fyne.App) Func {
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

func Link(w fyne.Window, input *widget.Entry) Func {
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
