package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xuyang-lee/wormhole/handle"
)

var takeCmd = &cobra.Command{
	Use:   "take",
	Short: "receive file or text from other hole",
	Long:  "start a background server to receive file or text",
	Args:  cobra.NoArgs,
	Run:   handle.Take,
}

func initTakeCmd() {
	takeCmd.Flags().StringVar(&handle.TakeHoleDir, "dir", "", "file will be save at dir")
	takeCmd.Flags().IntVarP(&handle.TakeHolePort, "port", "p", 0, "listen to port")
}
