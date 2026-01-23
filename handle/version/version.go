package version

import (
	"fmt"
	"github.com/xuyang-lee/wormhole/config/env"
)

func Show() {
	fmt.Println(env.ServerName, env.Version)
}
