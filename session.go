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
	// global session storage controller
	globalStore Storage
	// global Configure controller
	globalConfig *config
	// concurrent safe mutex
	mux sync.Mutex
	// Universal error message
	ErrKeyNoData      = errors.New("key no data")
	ErrSessionNoData  = errors.New("session no data")
	ErrIsEmpty        = errors.New("key OR session id is empty")
	ErrAlreadyExpired = errors.New("session already expired")
)

// Values is session item value
type Values map[string]interface{}

// Session is web session struct
type Session struct {
	Values
	ID         string
	CreateTime time.Time
	ExpireTime time.Time
}

// GetSession Get session data from the Request
func GetSession(w http.ResponseWriter, req *http.Request) (*Session, error) {
	debug.trace("Request:", req)

	mux.Lock()
	defer mux.Unlock()

	var session Session
	cookie, err := req.Cookie(globalConfig.CookieName)
	if cookie == nil || err != nil {
		debug.trace("cookie is empty:", cookie)
		return createSession(w, cookie, &session)
	}

	if len(cookie.Value) >= 73 {
		session.ID = cookie.Value
		if globalStore.Read(&session) != nil {
			return createSession(w, cookie, &session)
		}
	}

	debug.trace("session:", session)
	return &session, nil
}

// Sync save data modify
func (s *Session) Sync() error {
	debug.trace("session sync:", s)
	return globalStore.Write(s)
}

func createSession(w http.ResponseWriter, cookie *http.Cookie, session *Session) (*Session, error) {
	debug.trace("begin create session", session)

	session = NewSession()
	if cookie == nil {
		cookie = factory()
	}
	cookie.Value = session.ID
	cookie.MaxAge = int(globalConfig.LifeTime) / 1e9
	if err := globalStore.Write(session); err != nil {
		return nil, err
	}

	debug.trace("cookie:", cookie)

	http.SetCookie(w, cookie)

	debug.trace("end create session", session)
	return session, nil
}

// factory return default config pointer
func factory() *http.Cookie {
	return &http.Cookie{
		Domain:   globalConfig.Domain,
		Path:     globalConfig.DomainPath,
		Name:     globalConfig.CookieName,
		Secure:   globalConfig.Secure,
		HttpOnly: globalConfig.HttpOnly,
	}
}

// genreate session uuid length 73
func uuid73() string {
	return fmt.Sprintf("%s-%s", uuid.New().String(), uuid.New().String())
}

// NewSession return new session
func NewSession() *Session {
	nowTime := time.Now()
	return &Session{
		ID:         uuid73(),
		Values:     make(Values),
		CreateTime: nowTime,
		ExpireTime: nowTime.Add(lifeTime),
	}
}

// Expired check current session whether expire
func (s *Session) Expired() bool {
	return time.Duration(s.ExpireTime.UnixNano()) <= time.Duration(time.Now().UnixNano())
}

// Open Initialize storage with custom configuration
func Open(opt Configure) {
	debug.trace("open:", opt)

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
