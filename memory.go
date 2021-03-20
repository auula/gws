package sessionx

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Memory struct {
	sync.Mutex
	sessions map[string]*Session
}

func (m *Memory) Reader(s *Session) ([]byte, error) {
	m.Lock()
	defer m.Unlock()

	if ele, ok := m.sessions[s.ID]; ok {
		return encoder(ele)
	}

	return nil, fmt.Errorf("id %s not session data", s.ID)
}

func (m *Memory) Create(s *Session) ([]byte, error) {
	m.Lock()
	defer m.Unlock()
	value := make(map[string]interface{}, 8)
	if m.sessions == nil {
		m.sessions = make(map[string]*Session, 512*runtime.NumCPU())
	}
	s.Data = value
	s.Expires = time.Now().Add(time.Minute * 30)
	m.sessions[s.ID] = s
	return encoder(s)
}

func (m *Memory) Delete(s *Session) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.sessions[s.ID]; ok {
		delete(m.sessions, s.ID)
		return nil
	}
	return fmt.Errorf("id %s not find session data", s.ID)
}

func (m *Memory) Remove(s *Session, key string) error {
	m.Lock()
	defer m.Unlock()
	if ele, ok := m.sessions[s.ID]; ok {
		delete(ele.Data, key)
		return nil
	}
	return fmt.Errorf("id %s not find session data", s.ID)
}

func (m *Memory) Update(s *Session) ([]byte, error) {
	m.Lock()
	defer m.Unlock()
	if ele, ok := m.sessions[s.ID]; ok {
		ele.Data = s.Data
		ele.Expires = time.Now().Add(time.Minute * 30)
		//m.sessions[s.ID] = ele
		return encoder(ele)
	}
	return nil, fmt.Errorf("id %s updated session fail", s.ID)
}
