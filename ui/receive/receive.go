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
				common.ScrollToBottom(msgVScroll)
			case websocket.BinaryMessage:
			case -1: // websocket协议，(收到websocket.CloseMessage时链接已关闭，接受方收到的是-1而不是websocket.CloseMessage)
				common.AddSystemMessage(messageList, "关闭链接！")
				common.ScrollToBottom(msgVScroll)
				break readCh
			default:
				common.AddSystemMessage(messageList, fmt.Sprintf("unknown type:%d msg: %s", metaData.MsgType, string(metaData.Body)))
				common.ScrollToBottom(msgVScroll)
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
			common.ScrollToBottom(msgVScroll)
		case websocket.BinaryMessage:
		case websocket.CloseMessage:
			common.AddSystemMessage(messageList, "对方主动关闭链接！")
			common.ScrollToBottom(msgVScroll)
			common.CurLinkKey = ""
		default:
			common.AddSystemMessage(messageList, fmt.Sprintf("unknown type msg: %s", string(metaData.Body)))
			common.ScrollToBottom(msgVScroll)
		}
	}
}
