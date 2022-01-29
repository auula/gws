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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	// Global session storage controller
	globalStore Storage
)

func init() {
	globalStore = nil
}

type Storage interface {
	Clean(sid string)
	Remove(sid string, key string) error
	Get(sid string, key string, obj interface{}) (err error)
	Save(sid string, key string, obj interface{}) (err error)
}

// Values is session item value
type Values map[string][]byte

type RamStore struct {
	mux   sync.Mutex
	store map[string]*Session
}

func (ram *RamStore) Save(sid string, key string, obj interface{}) (err error) {

	if err := isEmpty(sid, key); err != nil {
		return err
	}

	var bytes []byte

	if bytes, err = json.Marshal(obj); err != nil {
		return err
	}

	ram.mux.Lock()
	ram.store[sid].Data[key] = bytes
	ram.mux.Unlock()

	return nil
}

func (ram *RamStore) Get(sid string, key string, obj interface{}) (err error) {

	if err := isEmpty(sid, key); err != nil {
		return err
	}

	var bytes []byte

	if bs, ok := ram.store[sid].Data[key]; !ok {
		// 如果是空这个bs 也是空并且返回了
		bytes = bs
		return errors.New("key no data")
	}

	return json.Unmarshal(bytes, obj)
}

func (ram *RamStore) Remove(sid string, key string) error {

	if err := isEmpty(sid, key); err != nil {
		return err
	}
	ram.mux.Lock()
	delete(ram.store[sid].Data, key)
	ram.mux.Unlock()
	return nil
}

func (ram *RamStore) Clean(sid string) {

	if sid == "" {
		return
	}

	ram.mux.Lock()
	delete(ram.store, sid)
	ram.mux.Unlock()
}

// gc is ram garbage collection.
func (ram *RamStore) gc() {
	for {
		// 30 minute garbage collection.
		time.Sleep(lifeTime)
		for _, session := range ram.store {
			if session.Expired() {
				ram.mux.Lock()
				delete(ram.store, session.UUID)
				ram.mux.Unlock()
			}
		}
	}
}

func isEmpty(sid string, key string) error {
	if key == "" || sid == "" {
		return errors.New("key OR session id is empty")
	}
	return nil
}

type Session struct {
	UUID string
	Data Values
	http.Cookie
	CreateTime time.Duration
	ExpireTime time.Duration
}

func (s Session) Save(key string, obj interface{}) (err error) {
	if s.Expired() {
		return errors.New("session already expired")
	}
	if err = globalStore.Save(s.UUID, key, obj); err != nil {
		return
	}
	// 重置会话生命周期
	s.renew()
	return
}

func (s Session) Get(key string, obj interface{}) (err error) {
	if s.Expired() {
		return errors.New("session already expired")
	}
	if err = globalStore.Get(s.UUID, key, obj); err != nil {
		return
	}
	// 重置会话生命周期
	s.renew()
	return
}

func (s Session) Remove(key string) error {
	if s.Expired() {
		return errors.New("session already expired")
	}
	if err := globalStore.Remove(s.UUID, key); err != nil {
		return err
	}
	// 重置会话生命周期
	s.renew()
	return nil
}

func (s Session) Clean() {
	globalStore.Clean(s.UUID)
}

func (s Session) refresh() {

}

func (s Session) Migrate() (*Session, error) {
	return nil, nil
}

func (s *Session) Expired() bool {
	return s.ExpireTime <= time.Duration(time.Now().UnixNano())
}

func (s *Session) renew() {
	nowTime := time.Duration(time.Now().UnixNano())
	s.ExpireTime = nowTime + lifeTime
	s.CreateTime = nowTime
}

func Handler(w http.ResponseWriter, req *http.Request) *Session {
	return nil
}

func UUID73() string {
	return fmt.Sprintf("%s-%s", uuid.New().String(), uuid.New().String())
}
