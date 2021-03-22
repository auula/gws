// MIT License

// Copyright (c) 2021 Jarvib Ding

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package sessionx

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type memoryStore struct {
	sync.Mutex
	sessions map[string]*Session
}

func (m *memoryStore) Reader(s *Session) ([]byte, error) {
	m.Lock()
	defer m.Unlock()

	if ele, ok := m.sessions[s.ID]; ok {
		return encoder(ele)
	}

	return nil, fmt.Errorf("id %s not session data", s.ID)
}

func (m *memoryStore) Create(s *Session) ([]byte, error) {
	m.Lock()
	defer m.Unlock()
	if m.sessions == nil {
		m.sessions = make(map[string]*Session, 512*runtime.NumCPU())
	}
	if s.Data == nil {
		s.Data = make(map[string]interface{}, 8)
	}
	m.sessions[s.ID] = s
	return encoder(s)
}

func (m *memoryStore) Delete(s *Session) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.sessions[s.ID]; ok {
		delete(m.sessions, s.ID)
		return nil
	}
	return fmt.Errorf("id %s not find session data", s.ID)
}

func (m *memoryStore) Remove(s *Session, key string) error {
	m.Lock()
	defer m.Unlock()
	if ele, ok := m.sessions[s.ID]; ok {
		delete(ele.Data, key)
		return nil
	}
	return fmt.Errorf("id %s not find session data", s.ID)
}

func (m *memoryStore) Update(s *Session) ([]byte, error) {
	m.Lock()
	defer m.Unlock()
	if ele, ok := m.sessions[s.ID]; ok {
		ele.Data = s.Data
		ele.Expires = time.Now().Add(mgr.cfg.SessionLifeTime)
		//m.sessions[s.ID] = ele
		return encoder(ele)
	}
	return nil, fmt.Errorf("id %s updated session fail", s.ID)
}

func (m *memoryStore) GC() {
	// recycle your trash every 10 minutes
	for {
		time.Sleep(time.Minute)
		m.Lock()
		for s, session := range m.sessions {
			if time.Now().UnixNano() >= session.Expires.UnixNano() {
				// log.Println("session-id: ", s, "expired.")
				delete(m.sessions, s)
			}
		}
		m.Unlock()
		// log.Println("gc running...")
	}

}
