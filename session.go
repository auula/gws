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

	if len(key) <= 0 {
		return nil, fmt.Errorf("key '%s' invalid", key)
	}
	bytes, err := mgr.store.Reader(s)
	if err != nil {
		return nil, err
	}
	_ = decoder(bytes, s)
	if ele, ok := s.Data[key]; ok {
		return ele, nil
	}
	return nil, fmt.Errorf("key '%s' does not exist", key)
}

func (s *Session) Set(key string, v interface{}) error {
	if len(key) <= 0 {
		return fmt.Errorf("key '%s' invalid", key)
	}
	mux.Lock()
	if s.Data == nil {
		s.Data = make(map[string]interface{}, 8)
	}
	s.Data[key] = v
	mux.Unlock()
	if _, err := mgr.store.Update(s); err != nil {
		return err
	}
	return nil
}

func (s *Session) Remove(key string) error {

	if err := mgr.store.Remove(s, key); err != nil {
		return err
	}
	return nil
}

func (s *Session) Clean() error {
	if err := mgr.store.Delete(s); err != nil {
		return err
	}
	return nil
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
