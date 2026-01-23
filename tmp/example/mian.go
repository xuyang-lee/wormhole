package main

//
//import (
//	"crypto/md5"
//	"encoding/binary"
//	"encoding/hex"
//	"flag"
//	"fmt"
//	"io"
//	"net"
//	"os"
//	"path/filepath"
//	"time"
//)
//
//const (
//	bufferSize = 32 * 1024 // 32KB buffer
//	port       = 8888
//)
//
//type FileInfo struct {
//	NameLen  uint32
//	Name     string
//	Size     int64
//	Checksum string
//}
//
//func main() {
//	mode := flag.String("mode", "send", "运行模式: send 或 receive")
//	filePath := flag.String("file", "", "要发送的文件路径 (发送模式)")
//	serverIP := flag.String("ip", "", "接收端IP地址 (发送模式)")
//	saveDir := flag.String("dir", ".", "接收文件保存目录 (接收模式)")
//	flag.Parse()
//
//	if *mode == "send" {
//		if *filePath == "" || *serverIP == "" {
//			fmt.Println("发送模式需要指定 -file 和 -ip 参数")
//			fmt.Println("示例: go run main.go -mode=send -file=/path/to/file -ip=192.168.1.100")
//			return
//		}
//		sendFile(*filePath, *serverIP)
//	} else if *mode == "receive" {
//		receiveFile(*saveDir)
//	} else {
//		fmt.Println("无效的模式,请使用 send 或 receive")
//	}
//}
//
//// 发送文件
//func sendFile(filePath, serverIP string) {
//	// 打开文件
//	file, err := os.Open(filePath)
//	if err != nil {
//		fmt.Printf("打开文件失败: %v\n", err)
//		return
//	}
//	defer file.Close()
//
//	// 获取文件信息
//	fileInfo, err := file.Stat()
//	if err != nil {
//		fmt.Printf("获取文件信息失败: %v\n", err)
//		return
//	}
//
//	// 计算文件MD5
//	fmt.Println("正在计算文件校验和...")
//	hash := md5.New()
//	if _, err := io.Copy(hash, file); err != nil {
//		fmt.Printf("计算校验和失败: %v\n", err)
//		return
//	}
//	checksum := hex.EncodeToString(hash.Sum(nil))
//	file.Seek(0, 0) // 重置文件指针
//
//	// 连接到接收端
//	addr := fmt.Sprintf("%s:%d", serverIP, port)
//	fmt.Printf("正在连接到 %s...\n", addr)
//	conn, err := net.Dial("tcp", addr)
//	if err != nil {
//		fmt.Printf("连接失败: %v\n", err)
//		return
//	}
//	defer conn.Close()
//
//	fmt.Println("连接成功,开始发送文件...")
//
//	// 发送文件名长度
//	fileName := filepath.Base(filePath)
//	nameLen := uint32(len(fileName))
//	if err := binary.Write(conn, binary.LittleEndian, nameLen); err != nil {
//		fmt.Printf("发送文件名长度失败: %v\n", err)
//		return
//	}
//
//	// 发送文件名
//	if _, err := conn.Write([]byte(fileName)); err != nil {
//		fmt.Printf("发送文件名失败: %v\n", err)
//		return
//	}
//
//	// 发送文件大小
//	if err := binary.Write(conn, binary.LittleEndian, fileInfo.Size()); err != nil {
//		fmt.Printf("发送文件大小失败: %v\n", err)
//		return
//	}
//
//	// 发送校验和长度和校验和
//	checksumLen := uint32(len(checksum))
//	if err := binary.Write(conn, binary.LittleEndian, checksumLen); err != nil {
//		fmt.Printf("发送校验和长度失败: %v\n", err)
//		return
//	}
//	if _, err := conn.Write([]byte(checksum)); err != nil {
//		fmt.Printf("发送校验和失败: %v\n", err)
//		return
//	}
//
//	// 发送文件内容
//	buffer := make([]byte, bufferSize)
//	var sent int64
//	startTime := time.Now()
//
//	for {
//		n, err := file.Read(buffer)
//		if err != nil && err != io.EOF {
//			fmt.Printf("读取文件失败: %v\n", err)
//			return
//		}
//		if n == 0 {
//			break
//		}
//
//		if _, err := conn.Write(buffer[:n]); err != nil {
//			fmt.Printf("发送数据失败: %v\n", err)
//			return
//		}
//
//		sent += int64(n)
//		progress := float64(sent) / float64(fileInfo.Size()) * 100
//		fmt.Printf("\r发送进度: %.2f%% (%d/%d bytes)", progress, sent, fileInfo.Size())
//	}
//
//	elapsed := time.Since(startTime)
//	speed := float64(sent) / elapsed.Seconds() / 1024 / 1024
//
//	fmt.Printf("\n文件发送完成!\n")
//	fmt.Printf("文件名: %s\n", fileName)
//	fmt.Printf("大小: %d bytes\n", fileInfo.Size())
//	fmt.Printf("MD5: %s\n", checksum)
//	fmt.Printf("耗时: %v\n", elapsed)
//	fmt.Printf("平均速度: %.2f MB/s\n", speed)
//}
//
//// 接收文件
//func receiveFile(saveDir string) {
//	// 创建保存目录
//	if err := os.MkdirAll(saveDir, 0755); err != nil {
//		fmt.Printf("创建目录失败: %v\n", err)
//		return
//	}
//
//	// 监听端口
//	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
//	if err != nil {
//		fmt.Printf("监听端口失败: %v\n", err)
//		return
//	}
//	defer listener.Close()
//
//	// 获取本机IP
//	addrs, _ := net.InterfaceAddrs()
//	fmt.Println("接收端已启动,等待连接...")
//	fmt.Println("本机IP地址:")
//	for _, addr := range addrs {
//		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
//			if ipnet.IP.To4() != nil {
//				fmt.Printf("  %s\n", ipnet.IP.String())
//			}
//		}
//	}
//	fmt.Printf("监听端口: %d\n\n", port)
//
//	for {
//		conn, err := listener.Accept()
//		if err != nil {
//			fmt.Printf("接受连接失败: %v\n", err)
//			continue
//		}
//
//		fmt.Printf("收到来自 %s 的连接\n", conn.RemoteAddr())
//		go handleConnection(conn, saveDir)
//	}
//}
//
//func handleConnection(conn net.Conn, saveDir string) {
//	defer conn.Close()
//
//	// 接收文件名长度
//	var nameLen uint32
//	if err := binary.Read(conn, binary.LittleEndian, &nameLen); err != nil {
//		fmt.Printf("接收文件名长度失败: %v\n", err)
//		return
//	}
//
//	// 接收文件名
//	nameBytes := make([]byte, nameLen)
//	if _, err := io.ReadFull(conn, nameBytes); err != nil {
//		fmt.Printf("接收文件名失败: %v\n", err)
//		return
//	}
//	fileName := string(nameBytes)
//
//	// 接收文件大小
//	var fileSize int64
//	if err := binary.Read(conn, binary.LittleEndian, &fileSize); err != nil {
//		fmt.Printf("接收文件大小失败: %v\n", err)
//		return
//	}
//
//	// 接收校验和
//	var checksumLen uint32
//	if err := binary.Read(conn, binary.LittleEndian, &checksumLen); err != nil {
//		fmt.Printf("接收校验和长度失败: %v\n", err)
//		return
//	}
//	checksumBytes := make([]byte, checksumLen)
//	if _, err := io.ReadFull(conn, checksumBytes); err != nil {
//		fmt.Printf("接收校验和失败: %v\n", err)
//		return
//	}
//	expectedChecksum := string(checksumBytes)
//
//	fmt.Printf("正在接收文件: %s (大小: %d bytes)\n", fileName, fileSize)
//
//	// 创建文件
//	filePath := filepath.Join(saveDir, fileName)
//	file, err := os.Create(filePath)
//	if err != nil {
//		fmt.Printf("创建文件失败: %v\n", err)
//		return
//	}
//	defer file.Close()
//
//	// 接收文件内容
//	buffer := make([]byte, bufferSize)
//	var received int64
//	hash := md5.New()
//	startTime := time.Now()
//
//	for received < fileSize {
//		n, err := conn.Read(buffer)
//		if err != nil {
//			fmt.Printf("接收数据失败: %v\n", err)
//			return
//		}
//
//		if _, err := file.Write(buffer[:n]); err != nil {
//			fmt.Printf("写入文件失败: %v\n", err)
//			return
//		}
//
//		hash.Write(buffer[:n])
//		received += int64(n)
//		progress := float64(received) / float64(fileSize) * 100
//		fmt.Printf("\r接收进度: %.2f%% (%d/%d bytes)", progress, received, fileSize)
//	}
//
//	elapsed := time.Since(startTime)
//	speed := float64(received) / elapsed.Seconds() / 1024 / 1024
//
//	// 验证校验和
//	actualChecksum := hex.EncodeToString(hash.Sum(nil))
//	fmt.Printf("\n\n文件接收完成!\n")
//	fmt.Printf("保存路径: %s\n", filePath)
//	fmt.Printf("文件大小: %d bytes\n", fileSize)
//	fmt.Printf("耗时: %v\n", elapsed)
//	fmt.Printf("平均速度: %.2f MB/s\n", speed)
//	fmt.Printf("MD5校验: ")
//	if actualChecksum == expectedChecksum {
//		fmt.Printf("✓ 通过 (%s)\n\n", actualChecksum)
//	} else {
//		fmt.Printf("✗ 失败\n")
//		fmt.Printf("  期望: %s\n", expectedChecksum)
//		fmt.Printf("  实际: %s\n\n", actualChecksum)
//	}
//}
