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

	globalConfig *config

	defaultCookie *http.Cookie

	// Universal error message
	ErrKeyNoData      = errors.New("key no data")
	ErrIsEmpty        = errors.New("key OR session id is empty")
	ErrAlreadyExpired = errors.New("session already expired")
	ErrSessionNoData  = errors.New("session no data")
)

type Storage interface {
	Read(s *Session) (err error)
	Write(s *Session) (err error)
	Create(s *Session) (err error)
	Remove(s *Session) (err error)
}

func Open(opt Configer) {
	globalConfig = opt.Parse()
	defaultCookie = factory()

	switch globalConfig.store {
	case ram:
		globalStore = &RamStore{
			mux:   sync.Mutex{},
			store: make(map[string]*Session),
		}
	case rds:
		globalStore = nil
	}
}

// Values is session item value
type Values map[string][]byte

type RamStore struct {
	mux   sync.Mutex
	store map[string]*Session
}

func (ram *RamStore) Create(s *Session) (err error) {

	if err := isEmpty(s.UUID); err != nil {
		return err
	}

	ram.mux.Lock()
	ram.store[s.UUID] = s
	ram.mux.Unlock()

	return nil
}

func (ram *RamStore) Read(s *Session) (err error) {

	if err := isEmpty(s.UUID); err != nil {
		return err
	}

	if session, ok := ram.store[s.UUID]; ok {
		s.Data = session.Data
		return nil
	}

	return ErrSessionNoData
}

func (ram *RamStore) Write(s *Session) (err error) {

	if err := isEmpty(s.UUID); err != nil {
		return err
	}

	ram.mux.Lock()
	if session, ok := ram.store[s.UUID]; ok {
		session.Data = s.Data
		return nil
	}
	ram.mux.Unlock()

	return nil
}

func (ram *RamStore) Remove(s *Session) (err error) {

	if err := isEmpty(s.UUID); err != nil {
		return err
	}

	ram.mux.Lock()
	delete(ram.store, s.UUID)
	ram.mux.Unlock()

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
				delete(ram.store, session.UUID)
				ram.mux.Unlock()
			}
		}
	}
}

func isEmpty(str string) error {
	if str == "" {
		return ErrIsEmpty
	}
	return nil
}

type Session struct {
	UUID string
	Data Values
	mux  sync.Mutex
	// 可有可无的字段 可以倒推出来
	CreateTime time.Duration
	ExpireTime time.Duration
}

func (s *Session) Save(key string, obj interface{}) (err error) {

	if s.Expired() {
		return ErrAlreadyExpired
	}

	var bytes []byte

	s.mux.Lock()
	if s.Data == nil {
		s.Data = make(Values)
	}

	if bytes, err = json.Marshal(obj); err != nil {
		return
	}
	s.Data[key] = bytes
	s.mux.Unlock()

	return globalStore.Write(s)
}

func (s *Session) Get(key string, obj interface{}) (err error) {
	if s.Expired() {
		return ErrAlreadyExpired
	}

	var (
		bytes []byte
		ok    bool
	)

	if err = globalStore.Read(s); err != nil {
		return err
	}

	if bytes, ok = s.Data[key]; !ok {
		// 如果是空这个bs 也是空并且返回了
		return ErrKeyNoData
	}

	return json.Unmarshal(bytes, obj)
}

func (s *Session) Remove(key string) error {
	if s.Expired() {
		return ErrAlreadyExpired
	}

	s.mux.Lock()
	delete(s.Data, key)
	s.mux.Unlock()

	return globalStore.Write(s)
}

func (s *Session) Clean() {
	globalStore.Remove(s)
}

func (s *Session) Migrate() (*Session, error) {
	return nil, nil
}

func (s *Session) Expired() bool {
	return s.ExpireTime <= time.Duration(time.Now().UnixNano())
}

func createSession(w http.ResponseWriter, cookie *http.Cookie, session *Session) *Session {

	// 初始化session参数
	nowTime := time.Duration(time.Now().UnixNano())

	session.UUID = uuid73()
	session.CreateTime = nowTime
	session.ExpireTime = nowTime + lifeTime

	return session
}

func factory() *http.Cookie {
	return &http.Cookie{
		Domain:   globalConfig.Domain,
		Path:     globalConfig.DomainPath,
		Name:     globalConfig.CookieName,
		Secure:   globalConfig.Secure,
		HttpOnly: globalConfig.HttpOnly,
	}
}

func GetSession(w http.ResponseWriter, req *http.Request) *Session {

	var session Session
	cookie, err := req.Cookie(globalConfig.CookieName)
	if cookie == nil || err != nil || cookie.Value == "" {
		return createSession(w, cookie, &session)
	}

	if len(cookie.Value) >= 73 {

		session.UUID = cookie.Value
		if err := globalStore.Read(&session); err != nil {
			return createSession(w, cookie, &session)
		}
		cookie := factory()
		cookie.Value = session.UUID
		cookie.MaxAge = session.ExpireTime
		http.SetCookie(w, cookie)
	}

	return &session
}

func uuid73() string {
	return fmt.Sprintf("%s-%s", uuid.New().String(), uuid.New().String())
}

func (s *Session) renew() {
	s.mux.Lock()
	defer s.mux.Unlock()
	nowTime := time.Duration(time.Now().UnixNano())
	s.ExpireTime = nowTime + lifeTime
	s.CreateTime = nowTime
}
