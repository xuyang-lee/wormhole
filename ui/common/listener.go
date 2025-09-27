package common

type Listener struct {
	f func()
}

func (l *Listener) Notify(key string) error {
	CurLinkKey = key
	l.f()
	return nil
}

func NewListener(f func()) *Listener {
	return &Listener{f: f}
}
