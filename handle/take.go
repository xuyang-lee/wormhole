package handle

import (
	"fmt"
	"github.com/spf13/cobra"
	envConf "github.com/xuyang-lee/wormhole/config/env"
	"github.com/xuyang-lee/wormhole/handle/env"
	"github.com/xuyang-lee/wormhole/handle/take"
	"github.com/xuyang-lee/wormhole/utils"
	"os"
)

var (
	TakeHoleDir  string
	TakeHolePort int
)

func Take(cmd *cobra.Command, args []string) {
	parseTakeParam()

	fmt.Println("当前进程: ", os.Getpid())

	if err := take.Listen(TakeHolePort, TakeHoleDir); err != nil {
		utils.ExitWithErr(err)
		return
	}

}

func parseTakeParam() {
	// 都有指定值，直接返回
	if TakeHoleDir != "" && TakeHolePort != 0 {
		return
	}

	// 有未指定的值，加载文件
	_ = env.LoadEnvFromFile()

	// dir 未指定
	if TakeHoleDir == "" {
		TakeHoleDir = utils.IfElse(envConf.HoleDir != "", envConf.HoleDir, envConf.DefaultHoleDir)
	}

	// port未指定
	if TakeHolePort == 0 {
		TakeHolePort = utils.IfElse(envConf.HolePort != 0, envConf.HolePort, envConf.DefaultHolePort)
	}
}
