package iface

import (
	"encoding/binary"
	"io"
)

// ============================================
// 数据源抽象: dataSource
// ============================================

// DataSource 数据源接口(文件、文本等)
type DataSource interface {
	// Read 读取数据
	Read(p []byte) (n int, err error)

	Write(p []byte) (n int, err error)
	// Size 获取数据大小
	Size() uint64
	// Type 获取数据类型
	Type() DataType
	// Metadata 获取元数据
	Metadata() map[string]string
	// Close 关闭资源
	Close() error
}

// DataType 数据类型
type DataType uint8

const (
	DataTypeFile DataType = 1
	DataTypeText DataType = 2
	DataTypeDir  DataType = 3
)

// ============================================
// 帧抽象: frames --> conn
// 		  frames
// ============================================

// FrameType 帧类型
type FrameType uint8

const (
	FrameTypeStd FrameType = 0x01
)

// Frame 接口
type Frame interface {
	Encode(w io.Writer) (n int64, err error)
	Decode(r io.Reader) (n int64, err error)
	Size() int64
	Payload() []byte

	Protocol() ProtocolFrame
}

type ProtocolFrame struct {
	HeaderSize int
	MaxPayload int64
	ByteOrder  binary.ByteOrder
}

// ============================================
// 传输抽象层
// ============================================

// Transfer 传输抽象接口
type Transfer interface {
	// ReadData 从数据源读取并转换为传输格式
	ReadData(source DataSource) error

	// WriteData 将传输结果转为数据源
	WriteData() (DataSource, error)

	// 字节流写入链接
	WriteTo(w io.Writer) (int64, error)

	// 从字节流中读数据
	ReadFrom(r io.Reader) (int64, error)

	// ToFrames 转换为帧序列
	Frames() []Frame
}
