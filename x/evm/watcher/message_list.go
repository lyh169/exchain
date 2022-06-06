package watcher

import (
	"container/list"
	"sync"
)

type idleWatchMessagesList struct {
	mutex    sync.Mutex
	idleList *list.List
}

func newIdleWatchMessageList() *idleWatchMessagesList {
	return &idleWatchMessagesList{
		idleList: list.New(),
	}
}

func (l *idleWatchMessagesList) getWatchMessages() []WatchMessage {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.idleList.Len() == 0 {
		return make([]WatchMessage, 0)
	}
	front := l.idleList.Remove(l.idleList.Front())

	return front.([]WatchMessage)[:0]
}

func (l *idleWatchMessagesList) putWatchMessages(msg []WatchMessage) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.idleList.PushBack(msg)
}
