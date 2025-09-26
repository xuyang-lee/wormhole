package hole

import (
	"errors"
	"github.com/xuyang-lee/wormhole/hole/session"
	"sync"
)

var once sync.Once
var linkLock sync.RWMutex
var linkMap map[string]*session.Session

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
	defer linkLock.Unlock()
	linkMap[link.Uuid()] = link
}

func UnRegisterLink(key string) {
	linkLock.Lock()
	defer linkLock.Unlock()
	delete(linkMap, key)
}
