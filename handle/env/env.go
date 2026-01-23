package env

import (
	"bufio"
	"github.com/xuyang-lee/wormhole/config/env"
	"github.com/xuyang-lee/wormhole/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func LoadEnvFromFile() error {

	path, err := utils.ConfigFilePath()
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		parseKeyValue(k, v)
	}
	return nil
}

func parseKeyValue(k, v string) {
	switch k {
	case env.KeyDstAddr:
		env.DstAddr = v
	case env.KeyDstIp:
		env.DstIp = v
	case env.KeyDstPort:
		env.DstPort, _ = strconv.Atoi(v)
	case env.KeyHolePort:
		env.HolePort, _ = strconv.Atoi(v)
	case env.KeyHoleDir:
		env.HoleDir = v
	}
}

func GetAllEnv() string {
	var sb strings.Builder
	sb.WriteString(env.KeyDstAddr + "=" + env.DstAddr + "\n")
	sb.WriteString(env.KeyDstIp + "=" + env.DstIp + "\n")
	sb.WriteString(env.KeyDstPort + "=" + strconv.Itoa(env.DstPort) + "\n")
	sb.WriteString(env.KeyHolePort + "=" + strconv.Itoa(env.HolePort) + "\n")
	sb.WriteString(env.KeyHoleDir + "=" + env.HoleDir + "\n")
	return sb.String()
}

func StoreEnvToFile() error {

	path, err := utils.ConfigFilePath()
	if err != nil {
		return err
	}

	conf := GetAllEnv()

	if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(conf)
	return err
}
