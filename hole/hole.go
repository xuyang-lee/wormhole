package hole

import (
	"errors"
	"github.com/xuyang-lee/wormhole/hole/session"
	"sync"
)

type Listener interface {
	Notify(string) error
}

var once sync.Once
var linkLock sync.RWMutex
var linkMap map[string]*session.Session
var listeners []Listener

func GetLink(uuid string) (*session.Session, bool) {
	linkLock.RLock()
	defer linkLock.RUnlock()
	link, ok := linkMap[uuid]
	return link, ok
}

func GetRandLink() (*session.Session, error) {
	linkLock.RLock()
	defer linkLock.RUnlock()
	for _, link := range linkMap {
		return link, nil
	}
	return nil, errors.New("no link")
}

func RegisterLink(link *session.Session) {
	linkLock.Lock()
	linkMap[link.Uuid()] = link
	linkLock.Unlock()
	for _, l := range listeners {
		_ = l.Notify(link.Uuid())
	}
}

func UnRegisterLink(key string) {
	linkLock.Lock()
	defer linkLock.Unlock()
	delete(linkMap, key)
}

func RegisterListener(l Listener) {
	listeners = append(listeners, l)
}
