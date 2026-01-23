package utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"github.com/xuyang-lee/wormhole/config/env"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func ConfigFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, env.ServerName, "env"), nil
}

func DownloadDirPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Downloads", env.ServerName), nil
}

func InitDownloadDir() error {
	dir, err := DownloadDirPath()
	if err != nil {
		return err
	}

	return os.MkdirAll(dir, os.ModePerm)
}

func LocalIP() (ips []string) {
	addrList, _ := net.InterfaceAddrs()
	for _, addr := range addrList {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return
}

func ExitWithErr(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func IfElse[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

func HeartBeat(d time.Duration) {
	log.Printf("heart beat logger started")
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for range ticker.C {
		log.Printf("in service...")
	}
}
