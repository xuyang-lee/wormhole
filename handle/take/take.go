package take

import (
	"errors"
	"fmt"
	"github.com/xuyang-lee/wormhole/handle/handshake"
	datasrc "github.com/xuyang-lee/wormhole/model/data_source"
	iface "github.com/xuyang-lee/wormhole/model/interface"
	"github.com/xuyang-lee/wormhole/model/transfer"
	"github.com/xuyang-lee/wormhole/utils"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

var (
	ErrListenPortFail = errors.New("listen to port fail")
)

func Listen(port int, dir string) error {

	// 创建保存目录
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("创建目录失败: %v\n", err)
		return err
	}

	// 监听端口
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("监听端口[%d]失败: %v\n", port, err)
		return ErrListenPortFail
	}
	defer listener.Close()

	// 获取本机IP
	fmt.Println("接收端已启动,等待连接...")
	fmt.Println("本机IP地址:")
	addrList := utils.LocalIP()
	for _, addr := range addrList {
		fmt.Printf("  %s\n", addr)
	}
	fmt.Printf("监听端口: %d\n\n", port)

	go utils.HeartBeat(5 * time.Minute)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("接受连接失败: %v\n", err)
			continue
		}

		fmt.Printf("收到来自 %s 的连接\n", conn.RemoteAddr())
		go handleConnection(conn, dir)
	}

}

func handleConnection(conn net.Conn, dir string) {
	defer conn.Close()

	if err := utils.TimeLimitConn(conn, 5*time.Second, handshake.Hi); err != nil {
		log.Printf("协议错误: %v\n", err)
		return
	}

	trans := transfer.New()
	if _, err := trans.ReadFrom(conn); err != nil && !errors.Is(err, io.EOF) {
		log.Println("transfer ReadFrom err", err)
		return
	}

	ds, err := trans.WriteData()
	if err != nil {
		log.Println("transfer WriteData err", err.Error())
		return
	}

	switch ds.Type() {
	case iface.DataTypeText:
		text := ds.(*datasrc.TextSource)
		fmt.Println(text.String())
	case iface.DataTypeFile:
		file := ds.(*datasrc.FileSource)
		fileName := filepath.Join(dir, file.Filename())
		fmt.Println(fileName)
		err = file.Commit(fileName)
		if err != nil {
			log.Println("file commit err:", err.Error())
			return
		}
	default:
		log.Println(datasrc.ErrDataSourceUnknownDataType)
		return
	}

	return
}
