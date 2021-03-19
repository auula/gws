package sessionx

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type Memory struct {
	lock     *sync.Mutex
	sessions map[string]*list.Element
}

type item struct {
	ID      string
	Lock    *sync.Mutex
	Expires time.Time
	Data    map[string]interface{}
}

func (m *Memory) Reader(id string) ([]byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if ele, ok := m.sessions[id]; ok {
		return Encoder(ele.Value.(*item))
	}

	return nil, errors.New("id not session data")
}

func (m *Memory) Create(id string) ([]byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	return nil, errors.New("id not session data")
}
