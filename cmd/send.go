package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xuyang-lee/wormhole/handle"
	"github.com/xuyang-lee/wormhole/model"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send file or text",
	Long:  "send a file or text to another hole",
	Args:  argsHook,
	Run:   handle.Send,
}

func initSendCmd() {
	sendCmd.Flags().StringVar(&handle.SendAddr, "addr", "", "address to send file to")
	sendCmd.Flags().StringVar(&handle.SendIp, "ip", "", "ip to send file to")
	sendCmd.Flags().IntVarP(&handle.SendPort, "port", "p", 0, "address to send file to")
	sendCmd.Flags().VarP(&handle.SendMode, "mode", "m", "transform mode [file|dir|text]. default use HOLE_PREFERRED_MODE if HOLE_PREFERRED_MODE had set, else use 'file'")

	sendCmd.MarkFlagsMutuallyExclusive("addr", "ip")
	sendCmd.MarkFlagsMutuallyExclusive("addr", "port")

}

func argsHook(cmd *cobra.Command, args []string) error {
	handle.ParseSendParam()
	switch handle.SendMode {
	case model.SendModeText:
		if len(args) < 1 {
			return fmt.Errorf("mode=text 时至少需要一个参数作为文本内容")
		}
	case model.SendModeNone, model.SendModeFile, model.SendModeDir:
		if len(args) != 1 {
			return fmt.Errorf("mode=file 时必须且只能提供一个参数作为文件路径")
		}
	default:
		return fmt.Errorf("未知的 mode: %s (允许值: text 或 file)", string(handle.SendMode))
	}
	return nil
}
