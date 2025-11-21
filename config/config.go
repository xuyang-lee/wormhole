package config

import (
	"fmt"
	"github.com/xuyang-lee/wormhole/utils"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

type AppConfig struct {
	Port int    `yaml:"port"`
	Addr string `yaml:"addr"`
}

var Conf AppConfig

func InitAppConfig() {

	var err error
	var data []byte

	dir, _ := os.Getwd()
	fmt.Println(dir)

	url := filepath.Join(utils.GetFilePath(), "config", "config.yaml")

	if data, err = os.ReadFile(url); err != nil {
		panic(fmt.Sprintf("无法读取配置文件: %s, 错误: %v", url, err))
	}

	if err = yaml.Unmarshal(data, &Conf); err != nil {
		panic(fmt.Sprintf("解析配置文件失败: %v", err))
	}

	log.Printf("加载配置文件 成功\n")
	return

}
