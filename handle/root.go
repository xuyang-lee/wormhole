package handle

import (
	"github.com/spf13/cobra"
	"github.com/xuyang-lee/wormhole/config/param"
	"github.com/xuyang-lee/wormhole/handle/version"
)

func Root(cmd *cobra.Command, args []string) {
	if param.ShowVersion {
		version.Show()
		return
	}
	if param.ShowHelp {
		_ = cmd.Help()
		return
	}
	return
}
