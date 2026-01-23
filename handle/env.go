package handle

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xuyang-lee/ezList/list"
	envConf "github.com/xuyang-lee/wormhole/config/env"
	"github.com/xuyang-lee/wormhole/handle/env"
	"github.com/xuyang-lee/wormhole/utils"
	"strings"
)

func Env(cmd *cobra.Command, args []string) {

	err := env.LoadEnvFromFile()
	if err != nil {
		utils.ExitWithErr(err)
		return
	}

	conf := env.GetAllEnv()

	defaultHoleDir, err := utils.DownloadDirPath()
	if err != nil {
		utils.ExitWithErr(err)
		return
	}

	ips := utils.LocalIP()
	ips = list.Extract(ips, func(s string) string { return "\t" + s })

	fmt.Println("====================You Can Set====================")
	fmt.Printf(conf)
	fmt.Println("====================Can Not Set====================")
	fmt.Printf("DEFAULT_HOLE_PORT=%d\n", envConf.DefaultHolePort)
	fmt.Printf("DEFAULT_HOLE_DIR=%s\n", defaultHoleDir)
	fmt.Printf("LOCAL_IPS:\n%s\n", strings.Join(ips, "\n"))

}
