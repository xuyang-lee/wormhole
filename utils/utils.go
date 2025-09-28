package utils

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/google/uuid"
	"net"
	"strings"
)

func GetUUID() string {
	clientID := uuid.New().String()
	return clientID
}

func GetUint64UUID() uint64 {
	var buf [8]byte
	_, _ = rand.Read(buf[:])
	return binary.BigEndian.Uint64(buf[:])
}

// 外连法：Dial 一个公网地址，读取本地出口 IP
func OutboundIP() (string, error) {
	// 使用 UDP 不会真的发包（除非必要），只用于获取本地 socket 的地址
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String() // "192.168.x.y:xxxxx"
	// 截取 ip 部分
	if i := strings.LastIndex(localAddr, ":"); i != -1 {
		return localAddr[:i], nil
	}
	return localAddr, nil
}
