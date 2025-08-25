package client

import (
	"fmt"
	"github.com/gorilla/websocket"
)

var Linked bool

func Dial(addr string) error {
	Linked = true
	fmt.Println(addr)
	return nil

	var conn *websocket.Conn
	var err error
	var data []byte

	if conn, _, err = websocket.DefaultDialer.Dial(addr, nil); err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer conn.Close()

	if err = conn.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
		fmt.Println(err.Error())
		return err
	}

	for {
		if _, data, err = conn.ReadMessage(); err != nil {
			// 意外的关闭错误，输出代码
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				fmt.Println(err.Error())
				return err
			}
			// 预期的关闭，输出关闭提示，优雅退出
			fmt.Println("server close connect")
			return nil
		}
		fmt.Println(string(data))

	}
}
