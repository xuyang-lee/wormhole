package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xuyang-lee/wormhole/handle"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set config",
	Long:  "after setting config, you can run hole without some flags",
	Args:  cobra.NoArgs,
	Run:   handle.Set,
}

func initSetCmd() {
	setCmd.Flags().StringVar(&handle.SetDstAddr, "addr", "", "the address where you send file to. this flag will reset HOLE_DST_IP and HOLE_DST_PORT \n e.g. 'xxx.xxx.xxx.xxx:port'")
	setCmd.Flags().StringVar(&handle.SetDstIp, "ip", "", "the ip where you send file to. this flag will reset HOLE_DST_ADDR")
	setCmd.Flags().IntVar(&handle.SetDstPort, "dport", 0, "the port where you send file to. this flag will reset HOLE_DST_ADDR")
	setCmd.Flags().StringVar(&handle.SetHoleDir, "dir", "", "file which you get from hole will be saved to this dir")
	setCmd.Flags().IntVar(&handle.SetHolePort, "port", 0, "port that your hole listens on")

	setCmd.MarkFlagsMutuallyExclusive("addr", "ip")
	setCmd.MarkFlagsMutuallyExclusive("addr", "dport")
}
