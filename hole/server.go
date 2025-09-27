package hole

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xuyang-lee/wormhole/config"
	"github.com/xuyang-lee/wormhole/hole/session"
	"github.com/xuyang-lee/wormhole/webtools"
	"log"
)

func Init() {
	go runServer()
}

func runServer() {
	config.InitAppConfig()

	InitLinkMap()

	// 最原生的升级手段
	//http.HandleFunc("/ws", wsHandler)
	//http.ListenAndServe("0.0.0.0:7777", nil)

	router := gin.Default()
	// 只能使用GET方法，其他方法会导致升级失败，原因是websocket协议规定，使用GET方法握手升级
	//本质的升级方式
	//router.GET("/ws", ginHandler)
	// 封装方法升级
	router.GET("/ws", webtools.Upgrade(keepConnect))
	log.Println("running")
	if err := router.Run(fmt.Sprintf(":%d", config.Conf.Port)); err != nil {
		panic(err.Error())
	}
}

// 使用封装的方法，升级http为webSocket
func keepConnect(conn *websocket.Conn) {

	link := session.NewSession(conn)
	RegisterLink(link)
	defer UnRegisterLink(link.Uuid())

	link.Wait()

	//var data []byte
	//var err error
	//
	//if _, data, err = conn.ReadMessage(); err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//
	//fmt.Println(string(data))
	//sendMsg := append(data, '!')
	//if err = conn.WriteMessage(websocket.TextMessage, sendMsg); err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	var err error
	// 主动断链接，优雅退出
	err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bey bey"))
	if err != nil {
		log.Println("Write close error:", err)
	}

	return
}
