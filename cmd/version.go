package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xuyang-lee/wormhole/handle"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version info",
	Long:  "version info",
	Args:  cobra.NoArgs,
	Run:   handle.Version,
}

func initVersionCmd() {

}
