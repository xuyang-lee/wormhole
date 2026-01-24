package datasrc

import (
	"bufio"
	"errors"
	iface "github.com/xuyang-lee/wormhole/model/interface"
	"io"
)

var (
	ErrDataSourceMismatch        = errors.New("data source mismatch")
	ErrDataSourceUnknownDataType = errors.New("unknown data type")
)

func GetFrameOfType(typ iface.DataType) iface.DataSource {
	switch typ {
	case iface.DataTypeFile:
		return &FileSource{}
	case iface.DataTypeText:
		return &TextSource{}
	default:
		return &FileSource{}
	}
}

func LoadToDataSource(r io.Reader) (iface.DataSource, error) {

	bufReader := bufio.NewReader(r)
	typ, err := bufReader.Peek(1)
	if err != nil {
		return nil, err
	}
	var data iface.DataSource
	switch iface.DataType(typ[0]) {
	case iface.DataTypeFile:
		data = NewEmptyFileSource()
	case iface.DataTypeText:
		data = NewEmptyTextSource()
	default:
		return nil, ErrDataSourceUnknownDataType
	}
	_, err = io.Copy(data, bufReader)
	if err != nil {
		return nil, err
	}
	return data, nil
}
