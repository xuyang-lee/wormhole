package hole

import (
	"errors"
	"github.com/xuyang-lee/wormhole/hole/session"
	"sync"
)

type linkController struct {
	sync.RWMutex
	once      sync.Once
	linkMap   map[string]*session.Session
	listeners []Listener
}

type Listener interface {
	Notify(string) error
}

var linkMap linkController

func InitLinkMap() {
	linkMap.initLinkMap()
}

func GetLink(uuid string) (*session.Session, bool) {
	return linkMap.getLink(uuid)
}

func GetRandLink() (*session.Session, error) {
	return linkMap.getRandLink()
}

func RegisterLink(link *session.Session) {
	linkMap.registerLink(link)
	return
}

func UnRegisterLink(key string) {
	linkMap.unRegisterLink(key)
	return
}

func RegisterListener(l Listener) {
	linkMap.registerListener(l)
}

func (linkCtrl *linkController) initLinkMap() {
	linkCtrl.once.Do(func() {
		linkCtrl.Lock()
		defer linkCtrl.Unlock()
		linkCtrl.linkMap = make(map[string]*session.Session, 10)
	})
}

func (linkCtrl *linkController) getLink(uuid string) (*session.Session, bool) {
	linkCtrl.RLock()
	defer linkCtrl.RUnlock()
	link, ok := linkCtrl.linkMap[uuid]
	return link, ok
}

func (linkCtrl *linkController) getRandLink() (*session.Session, error) {
	linkCtrl.RLock()
	defer linkCtrl.RUnlock()
	for _, link := range linkCtrl.linkMap {
		return link, nil
	}
	return nil, errors.New("no link")
}

func (linkCtrl *linkController) registerLink(link *session.Session) {
	linkCtrl.Lock()
	linkCtrl.linkMap[link.Uuid()] = link
	linkCtrl.Unlock()
	for _, l := range linkCtrl.listeners {
		_ = l.Notify(link.Uuid())
	}
}

func (linkCtrl *linkController) unRegisterLink(key string) {
	linkCtrl.Lock()
	defer linkCtrl.Unlock()
	delete(linkCtrl.linkMap, key)
}

func (linkCtrl *linkController) registerListener(l Listener) {
	linkCtrl.listeners = append(linkCtrl.listeners, l)
}
