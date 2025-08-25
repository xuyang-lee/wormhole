package webtools

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var u = websocket.Upgrader{
	HandshakeTimeout:  0,
	ReadBufferSize:    0,
	WriteBufferSize:   0,
	WriteBufferPool:   nil,
	Subprotocols:      nil,
	Error:             nil,
	CheckOrigin:       func(r *http.Request) bool { return true },
	EnableCompression: false,
}

func Upgrade(f func(conn *websocket.Conn), opts ...UpgradeOpt) gin.HandlerFunc {
	return func(c *gin.Context) {

		var up = u
		for _, opt := range opts {
			opt(&up)
		}

		conn, err := up.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("[Upgrade] err: %v\n", err.Error())
			return
		}
		defer conn.Close()

		f(conn)
	}
}

type A struct {
}

func (a A) Header() http.Header {
	//TODO implement me
	return http.Header{}
}

func (a A) Write(bytes []byte) (int, error) {
	return len(bytes), nil
}

func (a A) WriteHeader(statusCode int) {
	return
}
