package handle

import (
	"github.com/spf13/cobra"
	envConf "github.com/xuyang-lee/wormhole/config/env"
	"github.com/xuyang-lee/wormhole/handle/env"
	"github.com/xuyang-lee/wormhole/utils"
)

var (
	SetDstAddr  string
	SetDstIp    string
	SetDstPort  int
	SetHoleDir  string
	SetHolePort int
)

func Set(cmd *cobra.Command, args []string) {
	if err := env.LoadEnvFromFile(); err != nil {
		utils.ExitWithErr(err)
		return
	}

	flags := cmd.Flags()

	if flags.Changed("addr") {
		envConf.DstAddr = SetDstAddr
		envConf.DstIp = ""
		envConf.DstPort = 0
	}
	if flags.Changed("ip") {
		envConf.DstIp = SetDstIp
		envConf.DstAddr = ""
	}
	if flags.Changed("dport") {
		envConf.DstPort = SetDstPort
		envConf.DstAddr = ""
	}
	if flags.Changed("dir") {
		envConf.HoleDir = SetHoleDir
	}
	if flags.Changed("port") {
		envConf.HolePort = SetHolePort
	}

	if err := env.StoreEnvToFile(); err != nil {
		utils.ExitWithErr(err)
		return
	}

}
