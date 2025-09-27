package common

type Listener struct {
}

func (l *Listener) Notify(key string) error {
	CurLinkKey = key
	return nil
}
