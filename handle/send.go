package handle

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	envConf "github.com/xuyang-lee/wormhole/config/env"
	"github.com/xuyang-lee/wormhole/handle/env"
	"github.com/xuyang-lee/wormhole/model"
	datasrc "github.com/xuyang-lee/wormhole/model/data_source"
	iface "github.com/xuyang-lee/wormhole/model/interface"
	"github.com/xuyang-lee/wormhole/model/transfer"
	"github.com/xuyang-lee/wormhole/utils"
	"log"
	"net"
)

var (
	SendAddr string
	SendIp   string
	SendPort int
	SendMode model.SendModeFlag = model.SendModeFile
)

func Send(cmd *cobra.Command, args []string) {

	SendAddr = getAddr()
	arg := args[0]

	// 获取数据源
	var ds iface.DataSource
	var err error
	switch SendMode {
	case model.SendModeFile:
		ds, err = datasrc.NewFileSource(arg)
	case model.SendModeText:
		ds = datasrc.NewTextSource(arg)
	case model.SendModeDir:
		err = errors.New("dir mode is not currently supported")
	default:
		ds, err = datasrc.NewFileSource(arg)
	}
	if err != nil {
		utils.ExitWithErr(err)
	}

	// 获取转换器，并读取数据源
	trans := transfer.New()
	if err = trans.ReadData(ds); err != nil {
		utils.ExitWithErr(err)
	}

	// 连接到接收端
	fmt.Printf("正在连接到 %s...\n", SendAddr)
	conn, err := net.Dial("tcp", SendAddr)
	if err != nil {
		log.Printf("连接失败: %v\n", err)
		return
	}
	defer conn.Close()

	// 发送数据
	if _, err = trans.WriteTo(conn); err != nil {
		log.Printf("发送数据错误 ...\n")
		utils.ExitWithErr(err)
		return
	}

	fmt.Printf("发送成功 ...\n")
	return
}

func getAddr() string {
	// 指定了addr
	if SendAddr != "" {
		return SendAddr
	}

	// 指定了ip或port,其中之一
	if SendIp != "" || SendPort != 0 {
		var ip string
		var port int
		// 都指定了
		if SendIp != "" && SendPort != 0 {
			return fmt.Sprintf("%s:%d", SendIp, SendPort)
		}

		//有未指定的，加载配置文件
		_ = env.LoadEnvFromFile()
		// 优先级 SendIp> 配置文件 DstIp
		ip = utils.IfElse(SendIp != "", SendIp, envConf.DstIp)
		// 优先级 SendPort > 配置文件 DstPort > 默认 DefaultHolePort
		port = utils.IfElse(SendPort != 0, SendPort, utils.IfElse(envConf.DstPort != 0, envConf.DstPort, envConf.DefaultHolePort))

		return fmt.Sprintf("%s:%d", ip, port)
	}

	// 啥都没指定,先加载配置文件
	_ = env.LoadEnvFromFile()
	// 先看有没有文件保存addr
	if envConf.DstAddr != "" {
		return envConf.DstAddr
	}

	//没有的话看看文件保存的ip和port

	// ip没得默认值，找到port直接拼接地址就行
	if envConf.DstPort != 0 {
		return fmt.Sprintf("%s:%d", envConf.DstIp, envConf.DstPort)
	}

	return fmt.Sprintf("%s:%d", envConf.DstIp, envConf.DefaultHolePort)
}
