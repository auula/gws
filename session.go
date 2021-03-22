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
	"sync"
	"time"
)

var (
	mux sync.Mutex
	mgr *manager
)

type Session struct {
	ID      string
	Expires time.Time
	Data    map[string]interface{}
}

func (s *Session) Get(key string) (interface{}, error) {
	err := mgr.store.Reader(s)
	if err != nil {
		return nil, err
	}
	if ele, ok := s.Data[key]; ok {
		return ele, nil
	}
	return nil, fmt.Errorf("key '%s' does not exist", key)
}

func (s *Session) Set(key string, v interface{}) error {
	mux.Lock()
	if s.Data == nil {
		s.Data = make(map[string]interface{}, 8)
	}
	s.Data[key] = v
	mux.Unlock()
	return mgr.store.Update(s)
}

func (s *Session) Remove(key string) error {
	return mgr.store.Remove(s, key)
}

func (s *Session) Clean() error {
	return mgr.store.Delete(s)
}

// async code

// type Session struct {
// 	ID      string
// 	Expires time.Time
// 	Data    map[string]interface{}
// 	// 错误一次性收集
// 	tasks []task
// 	collgroup.Group
// }

// func (s *Session) Get(key string) interface{} {
// 	s.tasks = append(s.tasks, func() error {
// 		bytes, err := mgr.store.Reader(s)
// 		if err != nil {
// 			return err
// 		}
// 		Decoder(bytes, s)
// 		return nil
// 	})
// 	if ele, ok := s.Data[key]; ok {
// 		return ele
// 	}
// 	return nil
// }
