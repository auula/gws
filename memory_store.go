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
	sessions sync.Map
}

func (m *memoryStore) Reader(s *Session) error {
	if ele, ok := m.sessions.Load(s.ID); ok {
		// bug 这个不能直接 s = ele 因为有map地址
		s.Data = ele.(*Session).Data
		return nil
	}
	// s = nil
	return fmt.Errorf("id `%s` not exist session data", s.ID)
}

func (m *memoryStore) Create(s *Session) error {
	m.sessions.Store(s.ID, s)
	return nil
}

func (m *memoryStore) Delete(s *Session) error {
	m.sessions.Delete(s.ID)
	return nil
}

func (m *memoryStore) Update(s *Session) error {
	if ele, ok := m.sessions.Load(s.ID); ok {
		// 为什么是交换data 因为我们不确定上层是否扩容换了地址
		ele.(*Session).Data = s.Data
		ele.(*Session).Expires = time.Now().Add(mgr.cfg.TimeOut)
		//m.sessions[s.ID] = ele
		return nil
	}
	return fmt.Errorf("id `%s` updated session fail", s.ID)
}

func (m *memoryStore) gc() {
	// recycle your trash every 10 minutes
	for {
		time.Sleep(time.Minute * 10)
		m.sessions.Range(func(key, value interface{}) bool {
			if time.Now().UnixNano() >= value.(*Session).Expires.UnixNano() {
				m.sessions.Delete(key)
			}
			return true
		})
		runtime.GC()
		// log.Println("gc running...")
	}

}
