package model

import "fmt"

// SendModeFlag 自定义枚举类型
type SendModeFlag string

const (
	SendModeFile SendModeFlag = "file"
	SendModeDir  SendModeFlag = "dir"
	SendModeText SendModeFlag = "text"
)

// String 实现 pflag.Value 接口
func (m *SendModeFlag) String() string {
	return string(*m)
}

// Set 实现 pflag.Value 接口 (设置值时会调用)
func (m *SendModeFlag) Set(value string) error {
	switch value {
	case "file", "dir", "text":
		*m = SendModeFlag(value)
		return nil
	default:
		return fmt.Errorf("无效的模式: %s (有效值: file, dir, text)", value)
	}
}

// Type 实现 pflag.Value 接口
func (m *SendModeFlag) Type() string {
	return "mode"
}
