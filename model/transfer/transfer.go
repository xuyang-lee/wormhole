package transfer

import (
	"errors"
	"github.com/xuyang-lee/wormhole/model/generator"
	"github.com/xuyang-lee/wormhole/model/interface"
	"io"
)

var _ iface.Transfer = (*StandardTransfer)(nil)

var (
	ErrTransferHasBeenUsed = errors.New("transfer has been used")
)

type StandardTransfer struct {
	generator *generator.StandardGenerator
	frames    []iface.Frame
}

func (st *StandardTransfer) ReadData(source iface.DataSource) error {
	frames, err := st.generator.ReadData(source)
	if err != nil {
		return err
	}
	st.frames = frames
	return nil
}

func (st *StandardTransfer) WriteData() (iface.DataSource, error) {
	return st.generator.WriteData(st.frames)
}

func (st *StandardTransfer) WriteTo(w io.Writer) (int64, error) {
	var (
		n   int64
		en  int64
		err error
	)
	for _, f := range st.frames {
		en, err = f.Encode(w)
		n += en
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (st *StandardTransfer) ReadFrom(r io.Reader) (n int64, err error) {

	if st.frames != nil {
		return 0, ErrTransferHasBeenUsed
	}

	err = st.generator.SetReader(r).Err()
	if err != nil {
		return 0, err
	}

	for st.generator.Next() {
		fr, err := st.generator.Frame()
		if err != nil {
			return st.generator.Offset(), err
		}
		st.frames = append(st.frames, fr)
	}

	if err = st.generator.Err(); err != nil {
		return st.generator.Offset(), err
	}
	return st.generator.Offset(), nil
}

func (st *StandardTransfer) Frames() []iface.Frame {
	return st.frames
}

func New() *StandardTransfer {
	return &StandardTransfer{
		generator: generator.NewGenerator(iface.FrameTypeStd),
	}
}
