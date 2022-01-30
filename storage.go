// MIT License

// Copyright (c) 2022 Leon Ding

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

package gws

import (
	"sync"
	"time"
)

type Storage interface {
	Read(s *Session) (err error)
	Write(s *Session) (err error)
	Create(s *Session) (err error)
	Remove(s *Session) (err error)
}

type RamStore struct {
	mux   sync.Mutex
	store map[string]*Session
}

func (ram *RamStore) Create(s *Session) (err error) {
	ram.mux.Lock()
	defer ram.mux.Unlock()
	ram.store[s.ID] = s
	return nil
}

func (ram *RamStore) Read(s *Session) (err error) {
	if session, ok := ram.store[s.ID]; ok {
		s.Values = session.Values
		return nil
	}
	return ErrSessionNoData
}

func (ram *RamStore) Write(s *Session) (err error) {
	ram.mux.Lock()
	defer ram.mux.Unlock()
	if session, ok := ram.store[s.ID]; ok {
		session.Values = s.Values
	}
	return nil
}

func (ram *RamStore) Remove(s *Session) (err error) {
	ram.mux.Lock()
	defer ram.mux.Unlock()
	delete(ram.store, s.ID)
	return nil
}

// gc is ram garbage collection.
func (ram *RamStore) gc() {
	for {
		// 30 / 2 minute garbage collection.
		// 这里可以并发优化 向消费通道里面发送
		time.Sleep(lifeTime / 2)
		for _, session := range ram.store {
			if session.Expired() {
				ram.mux.Lock()
				delete(ram.store, session.ID)
				ram.mux.Unlock()
			}
		}
	}
}
