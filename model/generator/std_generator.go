package generator

import (
	"bytes"
	"errors"
	datasrc "github.com/xuyang-lee/wormhole/model/data_source"
	"github.com/xuyang-lee/wormhole/model/frame"
	iface "github.com/xuyang-lee/wormhole/model/interface"
	"io"
)

var (
	ErrNeedNextFrame        = errors.New("need to call Next() before Frame()")
	ErrUnknownFrameType     = errors.New("unknown frame type")
	ErrNoFrames             = errors.New("no frames")
	ErrByteReaderAlreadySet = errors.New("byte reader has already set")
)

type OptFunc func(*StandardGenerator)
type StandardGenerator struct {
	typ iface.FrameType

	// next fields are used in 'byte stream -> frame' scene
	nextFrame iface.Frame
	r         io.Reader
	offset    int64
	err       error
}

func (fg *StandardGenerator) ReadData(src iface.DataSource) ([]iface.Frame, error) {

	switch fg.typ {
	case iface.FrameTypeStd:
		return frame.GenerateStandardFrames(src)
	default:
		return nil, ErrUnknownFrameType
	}

}

func (fg *StandardGenerator) WriteData(frames []iface.Frame) (iface.DataSource, error) {
	if len(frames) == 0 {
		return nil, ErrNoFrames
	}
	var readers []io.Reader
	for _, oneFrame := range frames {
		readers = append(readers, bytes.NewReader(oneFrame.Payload()))
	}
	r := io.MultiReader(readers...)

	dataSource, err := datasrc.LoadToDataSource(r)
	if err != nil {
		return nil, err
	}

	return dataSource, nil
}

func (fg *StandardGenerator) Next() bool {
	if fg.err != nil {
		return false
	}

	f := frame.GetFrameOfType(fg.typ)

	n, err := f.Decode(fg.r)
	fg.offset += n
	if err != nil {
		fg.err = err
		return false
	}

	fg.nextFrame = f
	return true
}

func (fg *StandardGenerator) Frame() (iface.Frame, error) {
	if fg.nextFrame != nil { // 已经存在合法帧，直接返回
		return fg.nextFrame, nil
	}

	// 下一帧为空，判断下原因

	// error 导致
	if fg.err != nil {
		return nil, fg.err
	}

	return nil, ErrNeedNextFrame
}

func NewGenerator(typ iface.FrameType, opts ...OptFunc) *StandardGenerator {
	g := &StandardGenerator{
		typ: typ,
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

func WithReader(r io.Reader) OptFunc {
	return func(g *StandardGenerator) {
		g.r = r
		return
	}
}

func (fg *StandardGenerator) SetReader(r io.Reader) *StandardGenerator {
	if fg.r == nil {
		fg.r = r
	} else {
		fg.err = ErrByteReaderAlreadySet
	}
	return fg
}

func (fg *StandardGenerator) Err() error {
	return fg.err
}

func (fg *StandardGenerator) Offset() int64 {
	return fg.offset
}
