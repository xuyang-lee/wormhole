package webtools

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type UpgradeOpt func(u *websocket.Upgrader)

func WithHandshakeTimeout(t time.Duration) UpgradeOpt {
	return func(u *websocket.Upgrader) {
		u.HandshakeTimeout = t
	}
}

func WithReadBufferSize(size int) UpgradeOpt {
	return func(u *websocket.Upgrader) {
		u.ReadBufferSize = size
	}
}

func WithWriteBufferSize(size int) UpgradeOpt {
	return func(u *websocket.Upgrader) {
		u.WriteBufferSize = size
	}
}

func WithWriteBufferPool(pool websocket.BufferPool) UpgradeOpt {
	return func(u *websocket.Upgrader) {
		u.WriteBufferPool = pool
	}
}

func WithSubprotocols(s []string) UpgradeOpt {
	return func(u *websocket.Upgrader) {
		u.Subprotocols = s
	}
}

func WithError(f func(http.ResponseWriter, *http.Request, int, error)) UpgradeOpt {
	return func(u *websocket.Upgrader) {
		u.Error = f
	}
}

func WithCheckOrigin(f func(*http.Request) bool) UpgradeOpt {
	return func(u *websocket.Upgrader) {
		u.CheckOrigin = f
	}
}

func WithEnableCompression(e bool) UpgradeOpt {
	return func(u *websocket.Upgrader) {
		u.EnableCompression = e
	}
}
