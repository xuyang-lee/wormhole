package frame

import (
	"errors"
	iface "github.com/xuyang-lee/wormhole/model/interface"
	"io"
)

func GetFrameOfType(typ iface.FrameType) iface.Frame {
	switch typ {
	case iface.FrameTypeStd:
		return &StandardFrame{}
	default:
		return &StandardFrame{}
	}
}

func GenerateStandardFrames(r io.Reader) ([]iface.Frame, error) {

	frames := make([]iface.Frame, 0)
	var seq uint32 = 0
	for {
		seq++
		buf := make([]byte, HeaderSizeStandardFrame)
		n, err := io.ReadFull(r, buf)

		//没有下一帧。正常结束
		if errors.Is(err, io.EOF) {
			break
		}

		// 尾帧，最后一点内容填不满整个缓冲了
		if errors.Is(err, io.ErrUnexpectedEOF) {
			f, _ := NewStandardFrame(buf[:n], seq)
			frames = append(frames, f)
			break
		}

		// 其他未知错误
		if err != nil {
			return frames, err
		}

		// 完美的一帧
		f, _ := NewStandardFrame(buf, seq)
		frames = append(frames, f)
	}

	return frames, nil

}
