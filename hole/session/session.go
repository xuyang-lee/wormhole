package session

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/xuyang-lee/wormhole/utils"
	"log"
	"sync"
	"sync/atomic"
)

type Session struct {
	//wg          sync.WaitGroup
	once        sync.Once
	closed      atomic.Bool
	ctx         context.Context
	cancel      context.CancelFunc
	uuid        string
	conn        *websocket.Conn
	sendChan    chan MateElem
	receiveChan chan MateElem
}

func NewSession(conn *websocket.Conn) *Session {

	ctx, cancel := context.WithCancel(context.Background())

	session := &Session{
		uuid:        utils.GetUUID(),
		conn:        conn,
		sendChan:    make(chan MateElem, 1024),
		receiveChan: make(chan MateElem, 1024),
		ctx:         ctx,
		cancel:      cancel,
	}
	session.closed.Store(false)

	go session.sender()
	go session.receiver()

	return session
}

func (session *Session) sender() {
	defer session.close()

	for {
		if session.closed.Load() {
			return
		}
		select {
		case sendElem, ok := <-session.sendChan:
			if !ok {
				return
			}
			if err := session.conn.WriteMessage(sendElem.MsgType, sendElem.Body); err != nil {
				// 意外的关闭错误，输出代码
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
					log.Println("[sender]: ", err.Error())
					return
				}
				// 预期的关闭，输出关闭提示，优雅退出
				log.Println("[sender]: ", "hole close connect")
				return
			}
		case <-session.ctx.Done():
			return
		}
	}
}

func (session *Session) receiver() {
	defer close(session.receiveChan)
	defer session.close()

	for {
		if session.closed.Load() {
			return
		}
		var msgType int
		var data []byte
		var err error
		if msgType, data, err = session.conn.ReadMessage(); err != nil {
			// 意外的关闭错误，输出代码
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				fmt.Println("[receiver]: ", err.Error())
				return
			}
			// 预期的关闭，输出关闭提示，优雅退出
			fmt.Println("[receiver]: ", "hole close connect")
			return
		}
		receiveElem := MateElem{MsgType: msgType, Body: data}
		select {
		case session.receiveChan <- receiveElem:
		case <-session.ctx.Done():
			return
		}

	}
}

func (session *Session) close() {
	session.once.Do(func() {
		session.closed.Store(true)
		session.cancel()
	})
}
