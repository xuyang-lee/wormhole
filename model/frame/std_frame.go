package frame

import (
	"bytes"
	"encoding/binary"
	"errors"
	iface "github.com/xuyang-lee/wormhole/model/interface"
	"hash/crc32"
	"io"
)

const (
	MagicStandardFrame      = 0x8016
	VarStandardFrame        = 0
	HeaderSizeStandardFrame = 17
	MaxPayloadStandardFrame = 1*1024*1024 - HeaderSizeStandardFrame // 最大负载大小: 1MB - 头部大小
)

const (
	_ = iota
	StdFrameTypeHead
	StdFrameTypeData // 实际除了Data都是吹牛逼的预留帧
	StdFrameTypeTail
	StdFrameTypePing
	StdFrameTypePong
)

const (
	StdFrameFlagOne = 1 << iota // 预留flags
	StdFrameFlagTwo
	StdFrameFlagThree
	StdFrameFlagFore
	StdFrameFlagFive
	StdFrameFlagSix
	StdFrameFlagSeven
	StdFrameFlagEight
)

var (
	ByteOrderStandardFrame binary.ByteOrder = binary.BigEndian
)

var (
	ErrFrameCrcMismatch     = errors.New("frame crc mismatch")
	ErrFramePayloadTooLarge = errors.New("frame payload too large")
)

// StandardFrame 标准帧实现
type StandardFrame struct {
	magic   uint16
	ver     uint8
	typ     uint8
	flags   uint8
	len     uint32
	seq     uint32
	crc     uint32
	payload []byte
}

func (sf *StandardFrame) Encode(w io.Writer) (n int64, err error) {
	if err = binary.Write(w, ByteOrderStandardFrame, sf.magic); err != nil {
		return
	}
	n += 2
	if err = binary.Write(w, ByteOrderStandardFrame, sf.ver); err != nil {
		return
	}
	n += 1
	if err = binary.Write(w, ByteOrderStandardFrame, sf.typ); err != nil {
		return
	}
	n += 1
	if err = binary.Write(w, ByteOrderStandardFrame, sf.flags); err != nil {
		return
	}
	n += 1
	if err = binary.Write(w, ByteOrderStandardFrame, sf.len); err != nil {
		return
	}
	n += 4
	if err = binary.Write(w, ByteOrderStandardFrame, sf.seq); err != nil {
		return
	}
	n += 4
	if err = binary.Write(w, ByteOrderStandardFrame, sf.crc); err != nil {
		return
	}
	n += 4
	wn, err := io.Copy(w, bytes.NewReader(sf.payload))
	n += wn
	if err != nil {
		return
	}
	return
}

func (sf *StandardFrame) Decode(r io.Reader) (n int64, err error) {
	if err = binary.Read(r, ByteOrderStandardFrame, &sf.magic); err != nil {
		return
	}
	n += 2
	if err = binary.Read(r, ByteOrderStandardFrame, &sf.ver); err != nil {
		return
	}
	n += 1
	if err = binary.Read(r, ByteOrderStandardFrame, &sf.typ); err != nil {
		return
	}
	n += 1
	if err = binary.Read(r, ByteOrderStandardFrame, &sf.flags); err != nil {
		return
	}
	n += 1
	if err = binary.Read(r, ByteOrderStandardFrame, &sf.len); err != nil {
		return
	}
	n += 4
	if err = binary.Read(r, ByteOrderStandardFrame, &sf.seq); err != nil {
		return
	}
	n += 4
	if err = binary.Read(r, ByteOrderStandardFrame, &sf.crc); err != nil {
		return
	}
	n += 4
	if sf.len > 0 {
		sf.payload = make([]byte, sf.len)
		var rn int
		rn, err = io.ReadFull(r, sf.payload)
		n += int64(rn)
		if err != nil {
			return
		}
	}

	if sf.crc != sf.calcCrc() {
		return n, ErrFrameCrcMismatch
	}

	return n, err
}

func (sf *StandardFrame) Size() int64 {
	return HeaderSizeStandardFrame + int64(sf.len)
}

func (sf *StandardFrame) Payload() []byte {
	buf := make([]byte, sf.len)
	copy(buf, sf.payload)
	return buf
}

func (sf *StandardFrame) Protocol() iface.ProtocolFrame {
	return iface.ProtocolFrame{
		HeaderSize: HeaderSizeStandardFrame,
		MaxPayload: MaxPayloadStandardFrame,
		ByteOrder:  ByteOrderStandardFrame,
	}
}

func (sf *StandardFrame) calcCrc() uint32 {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, ByteOrderStandardFrame, sf.magic)
	_ = binary.Write(buf, ByteOrderStandardFrame, sf.ver)
	_ = binary.Write(buf, ByteOrderStandardFrame, sf.typ)
	_ = binary.Write(buf, ByteOrderStandardFrame, sf.flags)
	_ = binary.Write(buf, ByteOrderStandardFrame, sf.len)
	_ = binary.Write(buf, ByteOrderStandardFrame, sf.seq)
	// without crc field
	_, _ = buf.Write(sf.payload)

	return crc32.ChecksumIEEE(buf.Bytes())
}

func NewStandardFrame(payload []byte, seq uint32) (*StandardFrame, error) {
	if len(payload) > MaxPayloadStandardFrame {
		return nil, ErrFramePayloadTooLarge
	}
	sf := &StandardFrame{
		magic:   MagicStandardFrame,
		ver:     VarStandardFrame,
		typ:     StdFrameTypeData,
		flags:   0,
		len:     uint32(len(payload)),
		seq:     seq,
		crc:     0,
		payload: payload,
	}
	sf.crc = sf.calcCrc()
	return sf, nil
}
