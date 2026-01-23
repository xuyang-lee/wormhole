package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xuyang-lee/wormhole/handle"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send file or text",
	Long:  "send a file or text to another hole",
	Args:  cobra.ExactArgs(1),
	Run:   handle.Send,
}

func initSendCmd() {
	sendCmd.Flags().StringVar(&handle.SendAddr, "addr", "", "address to send file to")
	sendCmd.Flags().StringVar(&handle.SendIp, "ip", "", "ip to send file to")
	sendCmd.Flags().IntVarP(&handle.SendPort, "port", "p", 0, "address to send file to")
	sendCmd.Flags().VarP(&handle.SendMode, "mode", "m", "transform mode [file|dir|text]. default is file")

	sendCmd.MarkFlagsMutuallyExclusive("addr", "ip")
	sendCmd.MarkFlagsMutuallyExclusive("addr", "port")

}

// ModeFlag 自定义枚举类型
type ModeFlag string

const (
	ModeFile ModeFlag = "file"
	ModeDir  ModeFlag = "dir"
	ModeText ModeFlag = "text"
)

// String 实现 pflag.Value 接口
func (m *ModeFlag) String() string {
	return string(*m)
}

// Set 实现 pflag.Value 接口 (设置值时会调用)
func (m *ModeFlag) Set(value string) error {
	switch value {
	case "file", "dir", "text":
		*m = ModeFlag(value)
		return nil
	default:
		return fmt.Errorf("无效的模式: %s (有效值: file, dir, text)", value)
	}
}

// Type 实现 pflag.Value 接口
func (m *ModeFlag) Type() string {
	return "mode"
}
