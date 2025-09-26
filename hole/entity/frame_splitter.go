package entity

import (
	"errors"
	"fmt"
	"github.com/xuyang-lee/wormhole/hole/protocol"
	"github.com/xuyang-lee/wormhole/utils"
	"io"
)

// FrameType
const (
	FrameTypeNormal = iota + 1
	FrameTypeEnd
)

// FrameVersion
const (
	FrameVersionV250926 = iota + 1
	FrameVersionVMAX
)

type FrameSplitter struct {
	bufSize int
}

func (fs *FrameSplitter) Split(r io.Reader, size int) (<-chan protocol.Frame, error) {
	if r == nil {
		return nil, fmt.Errorf("reader is nil")
	}

	if size < 0 {
		return nil, fmt.Errorf("invalid frame size: %d", size)
	}

	frameCh := make(chan protocol.Frame, fs.bufSize)

	go fs.split(r, size, frameCh)

	return frameCh, nil
}

func (fs *FrameSplitter) Combine(frames <-chan protocol.Frame, w io.Writer) error {
	if w == nil {
		return fmt.Errorf("writer is nil")
	}

	if frames == nil {
		return fmt.Errorf("frame channel is nil")
	}

	for frame := range frames {
		if frame == nil {
			continue // 跳过空帧
		}

		// 将 Frame 序列化为字节
		data, err := frame.Dump()
		if err != nil {
			return fmt.Errorf("frame dump error: %w", err)
		}

		// 写入数据
		if _, err := w.Write(data); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
	}

	return nil
}

func (fs *FrameSplitter) split(r io.Reader, size int, fch chan protocol.Frame) {
	defer close(fch)

	var index uint64
	uuid := utils.GetUint64UUID()
	var version uint16 = FrameVersionV250926

	buffer := make([]byte, size)
	for {
		index = index + 1
		// 读取完整的一帧数据
		n, err := io.ReadFull(r, buffer)
		if err != nil {
			if err == io.EOF {
				endFrame := NewFrame(uuid, index, FrameTypeEnd, version)
				// 发送尾帧，并结束
				fch <- endFrame
				break // 正常结束
			}
			if errors.Is(err, io.ErrUnexpectedEOF) {
				// 最后一帧可能不完整，但仍然处理
				endFrame := NewFrame(uuid, index, FrameTypeEnd, version)
				if n > 0 {
					//发送尾帧
					endFrame.SetBody(buffer[:n])
					fch <- endFrame
				}
				break
			}
			// 其他错误，记录并退出
			fmt.Printf("Read error: %v\n", err)
			break
		}

		// 成功读取完整一帧，解析为 Frame
		f := NewFrame(uuid, index, FrameTypeNormal, version)
		f.SetBody(buffer[:n])
		// 发送到通道
		fch <- f
	}

}

func (fs *FrameSplitter) WithBufSize(bs int) *FrameSplitter {
	if fs == nil {
		return fs
	}
	fs.bufSize = bs
	return fs
}

func NewFrameSplitter() *FrameSplitter {
	return &FrameSplitter{
		bufSize: 4096,
	}
}
