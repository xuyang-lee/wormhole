package datasrc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/xuyang-lee/wormhole/config/env"
	"github.com/xuyang-lee/wormhole/model/interface"
	"io"
	"log"
	"os"
)

var _ iface.DataSource = (*FileSource)(nil)

var (
	ErrFileAborted       = errors.New("file source already aborted")
	ErrNoTmpFileToCommit = errors.New("no temp file to commit")
)

// FileSource 文件数据源
//
// read结果为 1b[type]+1b[nSize]+8b[size]+filename+fileBytes
// write也按此格式实现
type FileSource struct {
	nSize    uint8
	size     uint64
	filename string
	file     *os.File
	err      error

	// for Read
	r io.Reader
	// for Write
	buf         []byte
	protoOffset uint64 // 已写入文件大小

	committed bool
	aborted   bool
}

// NewFileSource 创建文件数据源
func NewFileSource(filepath string) (*FileSource, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	f := &FileSource{
		nSize:    uint8(len(info.Name())),
		filename: info.Name(),
		size:     uint64(info.Size()),
		file:     file,
	}

	return f, nil
}

func (f *FileSource) Read(p []byte) (int, error) {
	if f.r == nil {
		f.generateReader()
	}
	return f.r.Read(p)
}

func (f *FileSource) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	n := 0
	// 写第一个字节，判断类型匹配
	if f.protoOffset == 0 {
		if iface.DataType(p[0]) != iface.DataTypeFile {
			return 0, ErrDataSourceMismatch
		} else {
			f.protoOffset++
			n++
		}
		f.buf = make([]byte, 0, 9) //初始化9byte存储nSize和size信息
	}

	// 填充buf
	for f.protoOffset < 10 && n < len(p) {
		f.buf = append(f.buf, p[n])
		f.protoOffset++
		n++
		if len(f.buf) == 9 {
			f.parseSize()
			break
		}
	}

	// 填充filename
	for f.protoOffset < uint64(10+f.nSize) && n < len(p) {
		f.buf = append(f.buf, p[n])
		f.protoOffset++
		n++
		if len(f.buf) == int(f.nSize) {
			f.parseFilename()
			break
		}
	}

	if n == len(p) {
		return n, nil
	}

	// 写入前检查文件是否正确
	if f.err != nil {
		return n, f.err
	}
	w, err := f.file.Write(p[n:])
	f.protoOffset = f.protoOffset + uint64(w)
	return n + w, err
}

func (f *FileSource) Size() uint64 {
	return f.size
}

func (f *FileSource) Type() iface.DataType {
	return iface.DataTypeFile
}

func (f *FileSource) Metadata() map[string]string {
	return map[string]string{
		"filename": f.filename,
		"size":     fmt.Sprintf("%d", f.size),
	}
}

func (f *FileSource) Close() error {
	return f.file.Close()
}

func (f *FileSource) generateReader() {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, f.Type())
	_ = binary.Write(buf, binary.BigEndian, f.nSize)
	_ = binary.Write(buf, binary.BigEndian, f.size)
	_, _ = io.Copy(buf, bytes.NewBufferString(f.filename))
	f.r = io.MultiReader(buf, f.file)
}

func (f *FileSource) parseSize() {
	if len(f.buf) >= 1 {
		f.nSize = f.buf[0]
	}
	if len(f.buf) >= 9 {
		f.size = binary.BigEndian.Uint64(f.buf[1:9])
		// size 解析完成后，可以创建临时文件了
		f.err = f.ensureTempFile()
	}
	f.buf = make([]byte, 0, f.nSize) // 初始化nSize byte存储filename 信息
}

func (f *FileSource) parseFilename() {
	if len(f.buf) >= int(f.nSize) {
		f.filename = string(f.buf[:f.nSize])
	}
}

func (f *FileSource) ensureTempFile() error {
	if f.file != nil {
		return nil
	}
	log.Println("create tmp file, dir is", env.HoleDir)
	tmp, err := os.CreateTemp(env.HoleDir, "file-*")
	if err != nil {
		return err
	}
	f.file = tmp
	return nil
}

func (f *FileSource) Commit(path string) error {
	if f.aborted {
		return ErrFileAborted
	}
	if f.committed {
		return nil
	}
	if f.file == nil {
		return ErrNoTmpFileToCommit
	}

	// 先关闭文件，确保数据刷到磁盘
	if err := f.file.Close(); err != nil {
		return err
	}

	// Rename 是原子操作（同一文件系统内）
	if err := os.Rename(f.file.Name(), path); err != nil {
		return err
	}

	f.committed = true
	f.file = nil
	return nil
}

func (f *FileSource) Abort() error {
	if f.aborted {
		return nil
	}
	f.aborted = true

	if f.file == nil {
		return nil
	}

	name := f.file.Name()

	// 先关闭
	_ = f.file.Close()

	// 再删除
	err := os.Remove(name)

	f.file = nil
	return err
}

func (f *FileSource) Filename() string {
	return f.filename
}
