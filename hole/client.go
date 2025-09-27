package hole

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/xuyang-lee/wormhole/hole/session"
)

func Dial(addr string) error {
	var conn *websocket.Conn
	var err error

	if conn, _, err = websocket.DefaultDialer.Dial(addr, nil); err != nil {
		fmt.Println(err.Error())
		return err
	}

	go func() {
		defer conn.Close()
		link := session.NewSession(conn)
		RegisterLink(link)
		defer UnRegisterLink(link.Uuid())
		link.Wait()
	}()

	return nil
}
