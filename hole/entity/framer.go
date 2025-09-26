package entity

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	frameBaseSize = 8 + 8 + 4 + 4 + 2
)

type Frame struct {
	uid     uint64
	idx     uint64
	typ     uint32
	len     uint32
	ver     uint16
	payload []byte
}

func (f *Frame) SetUUID(u uint64) {
	if f != nil {
		f.uid = u
	}
}

func (f *Frame) SetIndex(u uint64) {
	if f != nil {
		f.idx = u
	}
}

func (f *Frame) SetType(u uint32) {
	if f != nil {
		f.typ = u
	}
}

func (f *Frame) SetVersion(u uint16) {
	if f != nil {
		f.ver = u
	}
}

func (f *Frame) SetBody(i []byte) {
	if f != nil {
		f.payload = make([]byte, len(i))
		copy(f.payload, i)
		f.len = uint32(len(i))
	}
}

func (f *Frame) UUID() uint64 {
	if f == nil {
		return 0
	}
	return f.uid
}

func (f *Frame) Index() uint64 {
	if f == nil {
		return 0
	}
	return f.idx
}

func (f *Frame) Type() uint32 {
	if f == nil {
		return 0
	}
	return f.typ
}

func (f *Frame) PayloadSize() uint32 {
	if f == nil {
		return 0
	}
	return f.len
}

func (f *Frame) Version() uint16 {
	if f == nil {
		return 0
	}
	return f.ver
}

func (f *Frame) HeadSize() int {
	return frameBaseSize
}

func (f *Frame) Payload() []byte {
	if f == nil {
		return nil
	}
	return f.payload
}

func (f *Frame) Dump() ([]byte, error) {
	if f == nil {
		return nil, errors.New("nil Frame")
	}

	totalSize := 8 + 4 + 8 + len(f.payload)
	buf := bytes.NewBuffer(make([]byte, 0, totalSize))

	if err := binary.Write(buf, binary.BigEndian, f.uid); err != nil {
		return nil, fmt.Errorf("dump uid err: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, f.idx); err != nil {
		return nil, fmt.Errorf("dump idx err: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, f.typ); err != nil {
		return nil, fmt.Errorf("dump typ err: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, f.len); err != nil {
		return nil, fmt.Errorf("dump len err: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, f.ver); err != nil {
		return nil, fmt.Errorf("dump ver err: %w", err)
	}

	if _, err := buf.Write(f.payload); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (f *Frame) Load(data []byte) error {
	if f == nil {
		return errors.New("nil Frame")
	}

	if len(data) < 20 { // 最小长度：idx(8) + typ(4) + uid(8) = 20字节
		return fmt.Errorf("data too short: %d bytes", len(data))
	}

	buf := bytes.NewReader(data)

	// 解析固定长度字段
	if err := binary.Read(buf, binary.BigEndian, &f.uid); err != nil {
		return fmt.Errorf("parse uid error: %w", err)
	}
	if err := binary.Read(buf, binary.BigEndian, &f.idx); err != nil {
		return fmt.Errorf("parse idx error: %w", err)
	}
	if err := binary.Read(buf, binary.BigEndian, &f.typ); err != nil {
		return fmt.Errorf("parse typ error: %w", err)
	}
	if err := binary.Read(buf, binary.BigEndian, &f.len); err != nil {
		return fmt.Errorf("parse len error: %w", err)
	}
	if err := binary.Read(buf, binary.BigEndian, &f.ver); err != nil {
		return fmt.Errorf("parse ver error: %w", err)
	}

	// 剩余的都是body数据
	f.payload = make([]byte, buf.Len())
	if _, err := buf.Read(f.payload); err != nil {
		return fmt.Errorf("parse payload error: %w", err)
	}

	return nil
}

func NewFrame(uid, idx uint64, typ uint32, ver uint16) *Frame {
	return &Frame{
		uid: uid,
		idx: idx,
		typ: typ,
		ver: ver,
	}
}

func NewEmptyFrame() *Frame {
	return &Frame{}
}
