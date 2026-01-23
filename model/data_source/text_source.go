package datasrc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/xuyang-lee/wormhole/model/interface"
	"io"
)

var _ iface.DataSource = (*TextSource)(nil)

// TextSource 文本数据源
type TextSource struct {
	text *bytes.Buffer
	size uint64

	// for Read
	r io.Reader
	// for Write
	buf         []byte
	protoOffset uint64 // 已写入文件大小
}

// NewTextSource 创建文本数据源
func NewTextSource(text string) *TextSource {
	return &TextSource{
		text: bytes.NewBuffer([]byte(text)),
		size: uint64(len(text)),
	}
}

func (t *TextSource) Read(p []byte) (int, error) {
	if t.r == nil {
		t.generateReader()
	}
	return t.r.Read(p)
}

func (t *TextSource) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	n := 0
	// 写第一个字节，判断类型匹配
	if t.protoOffset == 0 {
		if iface.DataType(p[0]) != iface.DataTypeText {
			return 0, ErrDataSourceMismatch
		} else {
			t.protoOffset++
			n++
		}
	}

	// 填充buf
	for t.protoOffset < 9 && n < len(p) {
		t.buf = append(t.buf, p[n])
		t.protoOffset++
		n++
		if len(t.buf) == 8 {
			t.parseSize()
			break
		}
	}

	if n == len(p) {
		return n, nil
	}

	w, err := t.text.Write(p[n:])
	t.protoOffset = t.protoOffset + uint64(w)
	return n + w, err
}

func (t *TextSource) Size() uint64 {
	return t.size
}

func (t *TextSource) Type() iface.DataType {
	return iface.DataTypeText
}

func (t *TextSource) Metadata() map[string]string {
	return map[string]string{
		"length": fmt.Sprintf("%d", t.text.Len()),
	}
}

func (t *TextSource) Close() error {
	return nil
}

func (t *TextSource) generateReader() {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, t.Type())
	_ = binary.Write(buf, binary.BigEndian, t.size)
	t.r = io.MultiReader(buf, t.text)
}

func (t *TextSource) parseSize() {
	if len(t.buf) >= 8 {
		t.size = binary.BigEndian.Uint64(t.buf[0:8])
	}
	t.buf = nil
}

func (t *TextSource) String() string {
	return string(t.text.Bytes())
}
