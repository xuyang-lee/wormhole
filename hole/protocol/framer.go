package protocol

import (
	"io"
)

type FrameSplitter interface {
	Split(r io.Reader, size int) (<-chan Frame, error)
	Combine(frames <-chan Frame, w io.Writer) error
}

type Frame interface {
	SetUUID(uint64)
	SetIndex(uint64)
	SetType(uint32)
	SetVersion(uint16)
	SetBody([]byte)

	Index() uint64
	Type() uint32
	UUID() uint64
	Version() uint16
	HeadSize() int
	PayloadSize() uint32
	Payload() []byte

	Dump() ([]byte, error)
	Load([]byte) error
}
