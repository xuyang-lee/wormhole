package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xuyang-lee/wormhole/webtools"
	"log"
	"net/http"
)

// 使用gin结合websocket，升级http的本质
func ginHandler(c *gin.Context) {
	wsHandler(c.Writer, c.Request)
}

// 最原生的方式，升级http到websocket
func wsHandler(w http.ResponseWriter, r *http.Request) {

	var data []byte
	var err error
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	if _, data, err = conn.ReadMessage(); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))
	sendMsg := append(data, '!')
	if err = conn.WriteMessage(websocket.TextMessage, sendMsg); err != nil {
		fmt.Println(err.Error())
		return
	}

	// 主动断链接，优雅退出
	err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bey bey"))
	if err != nil {
		log.Println("Write close error:", err)
	}

	return
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {

	// 最原生的升级手段
	//http.HandleFunc("/ws", wsHandler)
	//http.ListenAndServe("0.0.0.0:7777", nil)

	router := gin.Default()
	// 只能使用GET方法，其他方法会导致升级失败，原因是websocket协议规定，使用GET方法握手升级
	//本质的升级方式
	//router.GET("/ws", ginHandler)
	// 封装方法升级
	router.GET("/ws", webtools.Upgrade(exampleForWebTools))
	fmt.Println("running")
	if err := router.Run("0.0.0.0:7777"); err != nil {
		panic(err.Error())
	}

}

// 使用封装的方法，升级http为webSocket
func exampleForWebTools(conn *websocket.Conn) {

	var data []byte
	var err error

	if _, data, err = conn.ReadMessage(); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))
	sendMsg := append(data, '!')
	if err = conn.WriteMessage(websocket.TextMessage, sendMsg); err != nil {
		fmt.Println(err.Error())
		return
	}

	// 主动断链接，优雅退出
	err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bey bey"))
	if err != nil {
		log.Println("Write close error:", err)
	}

	return
}
