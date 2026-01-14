package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"
)

package main

import (
"crypto/md5"
"encoding/binary"
"encoding/hex"
"fmt"
"io"
"net"
"os"
"path/filepath"
"time"

"github.com/spf13/cobra"
"github.com/spf13/viper"
)

const (
	bufferSize = 32 * 1024
	version    = "1.0.0"
	configFile = ".worm.yaml"
)

var (
	cfgFile string
	ip      string
	port    int
)

func main2() {
	Execute()
}

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "worm",
	Short: "局域网文件传输工具",
	Long:  `Worm 是一个快速、安全的局域网文件传输工具`,
}

// versionCmd 版本命令
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("worm version %s\n", version)
	},
}

// envCmd 环境配置命令
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "显示当前配置",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("当前配置:")
		fmt.Printf("  默认IP: %s\n", viper.GetString("ip"))
		fmt.Printf("  默认端口: %d\n", viper.GetInt("port"))
		fmt.Printf("  接收目录: %s\n", viper.GetString("receive_dir"))
		fmt.Printf("  配置文件: %s\n", viper.ConfigFileUsed())
	},
}

// setCmd 设置配置命令
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "设置配置项",
	Long:  `设置默认配置,如: worm set --ip=192.168.1.1 --port=8080`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取标志值
		ip, _ := cmd.Flags().GetString("ip")
		port, _ := cmd.Flags().GetInt("port")
		receiveDir, _ := cmd.Flags().GetString("receive-dir")

		// 更新配置
		if ip != "" {
			viper.Set("ip", ip)
			fmt.Printf("✓ 设置默认IP: %s\n", ip)
		}
		if cmd.Flags().Changed("port") {
			viper.Set("port", port)
			fmt.Printf("✓ 设置默认端口: %d\n", port)
		}
		if receiveDir != "" {
			viper.Set("receive_dir", receiveDir)
			fmt.Printf("✓ 设置接收目录: %s\n", receiveDir)
		}

		// 保存配置
		home, _ := os.UserHomeDir()
		configPath := filepath.Join(home, configFile)
		if err := viper.WriteConfigAs(configPath); err != nil {
			fmt.Printf("保存配置失败: %v\n", err)
			return
		}
		fmt.Printf("✓ 配置已保存到: %s\n", configPath)
	},
}

// sendCmd 发送文件命令
var sendCmd = &cobra.Command{
	Use:   "send [文件路径]",
	Short: "发送文件",
	Long:  `发送文件到指定IP地址,如: worm send --ip=192.168.1.1 --port=8080 file.txt`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		// 获取IP和端口(优先使用命令行参数,其次使用配置文件)
		ip, _ := cmd.Flags().GetString("ip")
		if ip == "" {
			ip = viper.GetString("ip")
		}

		port, _ := cmd.Flags().GetInt("port")
		if !cmd.Flags().Changed("port") {
			port = viper.GetInt("port")
		}

		if ip == "" {
			fmt.Println("错误: 未指定目标IP地址")
			fmt.Println("请使用 --ip 参数或先运行 worm set --ip=<IP地址>")
			return
		}

		sendFile(filePath, ip, port)
	},
}

// receiveCmd 接收文件命令
var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "接收文件",
	Long:  `启动接收模式,等待文件传输,如: worm receive --port=8080 --dir=./downloads`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		if !cmd.Flags().Changed("port") {
			port = viper.GetInt("port")
		}

		dir, _ := cmd.Flags().GetString("dir")
		if dir == "" {
			dir = viper.GetString("receive_dir")
		}

		receiveFile(port, dir)
	},
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// 全局标志
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径 (默认: $HOME/.worm.yaml)")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "显示版本信息")

	// version 命令可以通过 -v 或 --version 触发
	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		v, _ := cmd.Flags().GetBool("version")
		if v {
			fmt.Printf("worm version %s\n", version)
			os.Exit(0)
		}
	}

	// set 命令的标志
	setCmd.Flags().String("ip", "", "设置默认IP地址")
	setCmd.Flags().Int("port", 0, "设置默认端口")
	setCmd.Flags().String("receive-dir", "", "设置默认接收目录")

	// send 命令的标志
	sendCmd.Flags().String("ip", "", "目标IP地址")
	sendCmd.Flags().IntP("port", "p", 0, "目标端口")

	// receive 命令的标志
	receiveCmd.Flags().IntP("port", "p", 0, "监听端口")
	receiveCmd.Flags().String("dir", "", "文件保存目录")

	// 添加子命令
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(receiveCmd)
}

// initConfig 初始化配置
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".worm")
	}

	// 设置默认值
	viper.SetDefault("port", 8080)
	viper.SetDefault("receive_dir", "./downloads")

	// 读取配置文件
	if err := viper.ReadInConfig(); err == nil {
		// 配置文件存在
	}
}

// sendFile 发送文件
func sendFile(filePath, serverIP string, port int) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("获取文件信息失败: %v\n", err)
		return
	}

	fmt.Println("正在计算文件校验和...")
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Printf("计算校验和失败: %v\n", err)
		return
	}
	checksum := hex.EncodeToString(hash.Sum(nil))
	file.Seek(0, 0)

	addr := fmt.Sprintf("%s:%d", serverIP, port)
	fmt.Printf("正在连接到 %s...\n", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("连接失败: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("连接成功,开始发送文件...")

	fileName := filepath.Base(filePath)
	nameLen := uint32(len(fileName))
	if err := binary.Write(conn, binary.LittleEndian, nameLen); err != nil {
		fmt.Printf("发送文件名长度失败: %v\n", err)
		return
	}

	if _, err := conn.Write([]byte(fileName)); err != nil {
		fmt.Printf("发送文件名失败: %v\n", err)
		return
	}

	if err := binary.Write(conn, binary.LittleEndian, fileInfo.Size()); err != nil {
		fmt.Printf("发送文件大小失败: %v\n", err)
		return
	}

	checksumLen := uint32(len(checksum))
	if err := binary.Write(conn, binary.LittleEndian, checksumLen); err != nil {
		fmt.Printf("发送校验和长度失败: %v\n", err)
		return
	}
	if _, err := conn.Write([]byte(checksum)); err != nil {
		fmt.Printf("发送校验和失败: %v\n", err)
		return
	}

	buffer := make([]byte, bufferSize)
	var sent int64
	startTime := time.Now()

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Printf("读取文件失败: %v\n", err)
			return
		}
		if n == 0 {
			break
		}

		if _, err := conn.Write(buffer[:n]); err != nil {
			fmt.Printf("发送数据失败: %v\n", err)
			return
		}

		sent += int64(n)
		progress := float64(sent) / float64(fileInfo.Size()) * 100
		fmt.Printf("\r发送进度: %.2f%% (%d/%d bytes)", progress, sent, fileInfo.Size())
	}

	elapsed := time.Since(startTime)
	speed := float64(sent) / elapsed.Seconds() / 1024 / 1024

	fmt.Printf("\n文件发送完成!\n")
	fmt.Printf("文件名: %s\n", fileName)
	fmt.Printf("大小: %d bytes\n", fileInfo.Size())
	fmt.Printf("MD5: %s\n", checksum)
	fmt.Printf("耗时: %v\n", elapsed)
	fmt.Printf("平均速度: %.2f MB/s\n", speed)
}

// receiveFile 接收文件
func receiveFile(port int, saveDir string) {
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		fmt.Printf("创建目录失败: %v\n", err)
		return
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("监听端口失败: %v\n", err)
		return
	}
	defer listener.Close()

	addrs, _ := net.InterfaceAddrs()
	fmt.Println("接收端已启动,等待连接...")
	fmt.Println("本机IP地址:")
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Printf("  %s\n", ipnet.IP.String())
			}
		}
	}
	fmt.Printf("监听端口: %d\n", port)
	fmt.Printf("保存目录: %s\n\n", saveDir)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("接受连接失败: %v\n", err)
			continue
		}

		fmt.Printf("收到来自 %s 的连接\n", conn.RemoteAddr())
		go handleConnection(conn, saveDir)
	}
}

func handleConnection(conn net.Conn, saveDir string) {
	defer conn.Close()

	var nameLen uint32
	if err := binary.Read(conn, binary.LittleEndian, &nameLen); err != nil {
		fmt.Printf("接收文件名长度失败: %v\n", err)
		return
	}

	nameBytes := make([]byte, nameLen)
	if _, err := io.ReadFull(conn, nameBytes); err != nil {
		fmt.Printf("接收文件名失败: %v\n", err)
		return
	}
	fileName := string(nameBytes)

	var fileSize int64
	if err := binary.Read(conn, binary.LittleEndian, &fileSize); err != nil {
		fmt.Printf("接收文件大小失败: %v\n", err)
		return
	}

	var checksumLen uint32
	if err := binary.Read(conn, binary.LittleEndian, &checksumLen); err != nil {
		fmt.Printf("接收校验和长度失败: %v\n", err)
		return
	}
	checksumBytes := make([]byte, checksumLen)
	if _, err := io.ReadFull(conn, checksumBytes); err != nil {
		fmt.Printf("接收校验和失败: %v\n", err)
		return
	}
	expectedChecksum := string(checksumBytes)

	fmt.Printf("正在接收文件: %s (大小: %d bytes)\n", fileName, fileSize)

	filePath := filepath.Join(saveDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("创建文件失败: %v\n", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, bufferSize)
	var received int64
	hash := md5.New()
	startTime := time.Now()

	for received < fileSize {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("接收数据失败: %v\n", err)
			return
		}

		if _, err := file.Write(buffer[:n]); err != nil {
			fmt.Printf("写入文件失败: %v\n", err)
			return
		}

		hash.Write(buffer[:n])
		received += int64(n)
		progress := float64(received) / float64(fileSize) * 100
		fmt.Printf("\r接收进度: %.2f%% (%d/%d bytes)", progress, received, fileSize)
	}

	elapsed := time.Since(startTime)
	speed := float64(received) / elapsed.Seconds() / 1024 / 1024

	actualChecksum := hex.EncodeToString(hash.Sum(nil))
	fmt.Printf("\n\n文件接收完成!\n")
	fmt.Printf("保存路径: %s\n", filePath)
	fmt.Printf("文件大小: %d bytes\n", fileSize)
	fmt.Printf("耗时: %v\n", elapsed)
	fmt.Printf("平均速度: %.2f MB/s\n", speed)
	fmt.Printf("MD5校验: ")
	if actualChecksum == expectedChecksum {
		fmt.Printf("✓ 通过 (%s)\n\n", actualChecksum)
	} else {
		fmt.Printf("✗ 失败\n")
		fmt.Printf("  期望: %s\n", expectedChecksum)
		fmt.Printf("  实际: %s\n\n", actualChecksum)
	}
}