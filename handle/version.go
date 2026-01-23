package handle

import (
	"github.com/spf13/cobra"
	"github.com/xuyang-lee/wormhole/handle/version"
)

func Version(cmd *cobra.Command, args []string) {
	version.Show()
}
