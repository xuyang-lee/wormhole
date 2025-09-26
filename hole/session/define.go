package session

type Sender interface {
	Send(int, []byte) error
}

type Receiver interface {
	Receive() *MateElem
}

type MateElem struct {
	MsgType int
	Body    []byte
}
