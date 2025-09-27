package receive

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/gorilla/websocket"
	"github.com/xuyang-lee/wormhole/hole"
	"github.com/xuyang-lee/wormhole/hole/session"
	"github.com/xuyang-lee/wormhole/ui/common"
	"time"
)

func Receive(app fyne.App, messageList *fyne.Container, msgVScroll *container.Scroll) {
	for {
		time.Sleep(10 * time.Second)
		link, ok := hole.GetLink(common.CurLinkKey)
		if !ok {
			time.Sleep(5 * time.Second)
			continue
		}

		ch := link.GetReceiveChannel()
	readCh:
		for metaData := range ch {
			switch metaData.MsgType {
			case websocket.TextMessage:
				common.AddReceiveMessageWithCopyButton(app, messageList, string(metaData.Body))
				msgVScroll.ScrollToBottom()
			case websocket.BinaryMessage:
			case websocket.CloseMessage:
				common.AddSystemMessage(messageList, "对方主动关闭链接！")
				msgVScroll.ScrollToBottom()
				break readCh
			default:
				common.AddSystemMessage(messageList, fmt.Sprintf("unknown type msg: %s", string(metaData.Body)))
				msgVScroll.ScrollToBottom()
			}
		}
	}
}

func TestReceive(app fyne.App, messageList *fyne.Container, msgVScroll *container.Scroll) {

	ch := make(chan session.MateElem)
	time.Sleep(10 * time.Second)
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)
			ch <- session.MateElem{
				MsgType: websocket.TextMessage,
				Body:    []byte(fmt.Sprintf("hello wormhole, this is the %dth msg", i+1)),
			}
		}
		ch <- session.MateElem{
			MsgType: websocket.CloseMessage,
			Body:    websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bey bey"),
		}
		close(ch)
	}()

	for metaData := range ch {
		switch metaData.MsgType {
		case websocket.TextMessage:
			common.AddReceiveMessageWithCopyButton(app, messageList, string(metaData.Body))
			msgVScroll.ScrollToBottom()
		case websocket.BinaryMessage:
		case websocket.CloseMessage:
			common.AddSystemMessage(messageList, "对方主动关闭链接！")
			msgVScroll.ScrollToBottom()
		default:
			common.AddSystemMessage(messageList, fmt.Sprintf("unknown type msg: %s", string(metaData.Body)))
			msgVScroll.ScrollToBottom()
		}
	}
}
