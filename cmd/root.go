package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xuyang-lee/wormhole/config/param"
	"github.com/xuyang-lee/wormhole/handle"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "hole",
	Short: "局域网文件传输工具",
	Long:  `Worm 是一个快速、安全的局域网文件传输工具`,
	Run:   handle.Root,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initRootCmd() {
	rootCmd.Flags().BoolVarP(&param.ShowVersion, "version", "v", false, "print hole version")
	//rootCmd.Flags().BoolVarP(&config.ShowHelp, "help", "h", false, "help")
	rootCmd.MarkFlagsMutuallyExclusive()

}

func init() {
	initVersionCmd()
	initSendCmd()
	initTakeCmd()
	initEnvCmd()
	initSetCmd()
	initRootCmd()

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(takeCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(setCmd)

}
