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

	mux sync.Mutex

	// Universal error message
	ErrKeyNoData      = errors.New("key no data")
	ErrSessionNoData  = errors.New("session no data")
	ErrIsEmpty        = errors.New("key OR session id is empty")
	ErrAlreadyExpired = errors.New("session already expired")
)

// Values is session item value
type Values map[string]interface{}

type Session struct {
	Values
	ID         string
	CreateTime time.Duration
	ExpireTime time.Duration
}

func GetSession(w http.ResponseWriter, req *http.Request) (*Session, error) {

	mux.Lock()
	defer mux.Unlock()

	var session Session
	cookie, err := req.Cookie(globalConfig.CookieName)
	if cookie == nil || err != nil {
		return createSession(w, cookie, &session)
	}

	if len(cookie.Value) >= 73 {
		session.ID = cookie.Value
		// 如果读不到说明存储里面没有
		if globalStore.Read(&session) != nil {
			return createSession(w, cookie, &session)
		}
	}
	return &session, nil
}

func (s *Session) Sync() error {
	return globalStore.Write(s)
}

func createSession(w http.ResponseWriter, cookie *http.Cookie, session *Session) (*Session, error) {

	// 初始化session参数
	nowTime := time.Duration(time.Now().UnixNano())
	session.ID = uuid73()
	session.CreateTime = nowTime
	session.ExpireTime = nowTime + lifeTime
	session.Values = make(Values)

	// 初始化cookie
	if cookie == nil {
		cookie = factory()
	}
	cookie.Value = session.ID
	cookie.MaxAge = int(globalConfig.LifeTime)
	http.SetCookie(w, cookie)
	if err := globalStore.Create(session); err != nil {
		return nil, err
	}

	return session, nil
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

func uuid73() string {
	return fmt.Sprintf("%s-%s", uuid.New().String(), uuid.New().String())
}

func (s *Session) Expired() bool {
	return s.ExpireTime <= time.Duration(time.Now().UnixNano())
}

func Open(opt Configure) {
	globalConfig = opt.Parse()
	switch globalConfig.store {
	case ram:
		globalStore = &RamStore{
			store: make(map[string]*Session),
			mux:   sync.Mutex{},
		}
	case rds:
		globalStore = nil
	}
}
