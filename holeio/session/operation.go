package session

import "errors"

func (session *Session) Send(msgType int, msg []byte) error {
	elem := MateElem{
		MsgType: msgType,
		Body:    msg,
	}
	select {
	case <-session.ctx.Done():
		return session.ctx.Err()
	case session.sendChan <- elem:
		return nil
	}

}

func (session *Session) Receive() *MateElem {
	msg := <-session.receiveChan
	return &msg
}

func (session *Session) TryReceive() (*MateElem, error) {
	select {
	case elem, ok := <-session.receiveChan:
		if !ok {
			return nil, errors.New("session closed")
		}
		return &elem, nil
	default:
		return nil, errors.New("no msg")
	}
}

func (session *Session) GetReceiveChannel() chan MateElem {
	return session.receiveChan
}

func (session *Session) Uuid() string {
	return session.uuid
}

func (session *Session) Wait() {
	<-session.ctx.Done()
}

func (session *Session) Close() {
	session.close()
}
