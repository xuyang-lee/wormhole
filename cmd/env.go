package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xuyang-lee/wormhole/handle"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "env config",
	Long:  "env config",
	Args:  cobra.NoArgs,
	Run:   handle.Env,
}

func initEnvCmd() {

}
