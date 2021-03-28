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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/barkimedes/go-deepcopy"
	"github.com/google/uuid"
)

type method func(f func())

var (
	rwm  sync.RWMutex
	mgr  *manager
	lock = map[string]method{
		"W": func(f func()) {
			rwm.Lock()
			defer rwm.Unlock()
			f()
		},
		"R": func(f func()) {
			rwm.RLock()
			defer rwm.RUnlock()
			f()
		},
	}
)

type Session struct {
	// sessionId
	ID string
	// session timeout
	Expires time.Time
	// map to store data
	Data map[string]interface{}
	_w   http.ResponseWriter
	// each session corresponds to a cookie
	Cookie *http.Cookie
}

// Get Retrieves the stored element data from the session via the key
func (s *Session) Get(key string) (interface{}, error) {
	err := mgr.store.Read(s)
	if err != nil {
		return nil, err
	}
	s.refreshCookie()
	if ele, ok := s.Data[key]; ok {
		return ele, nil
	}
	return nil, fmt.Errorf("key '%s' does not exist", key)
}

// Set Stores information in the session
func (s *Session) Set(key string, v interface{}) error {

	lock["W"](func() {
		if s.Data == nil {
			s.Data = make(map[string]interface{}, 8)
		}
		s.Data[key] = v
	})

	s.refreshCookie()
	return mgr.store.Update(s)
}

// Remove an element stored in the session
func (s *Session) Remove(key string) error {
	s.refreshCookie()

	lock["W"](func() {
		delete(s.Data, key)
	})

	return mgr.store.Update(s)
}

// Clean up all data for this session
func (s *Session) Clean() error {
	s.refreshCookie()
	return mgr.store.Remove(s)
}

// Handler Get session data from the Request
func Handler(w http.ResponseWriter, req *http.Request) *Session {

	// Take the session from the request
	var session Session
	session._w = w
	cookie, err := req.Cookie(mgr.cfg._cookie.Name)
	if err != nil || cookie == nil || len(cookie.Value) <= 0 {
		return createSession(w, cookie, &session)
	}

	// The length of the ID after being encoded is 73 bits
	if len(cookie.Value) >= 73 {
		session.ID = cookie.Value
		if mgr.store.Read(&session) != nil {
			return createSession(w, cookie, &session)
		}

		session.copy(mgr.cfg._cookie)
		session.Cookie.Value = session.ID
		session.Cookie.Expires = session.Expires
		session.refreshCookie()
	}

	return &session
}

func createSession(w http.ResponseWriter, cookie *http.Cookie, session *Session) *Session {
	// init session parameter
	session.ID = generateUUID()
	session.Expires = time.Now().Add(mgr.cfg.TimeOut)

	// 重置配置cookie模板
	session.copy(mgr.cfg._cookie)
	session.Cookie.Value = session.ID
	session.Cookie.Expires = session.Expires

	mgr.store.Create(session)

	session.refreshCookie()
	return session
}

// refreshCookie: Refresh the cookie session as long as there is an operation to reset the session life cycle
func (s *Session) refreshCookie() {
	s.Expires = time.Now().Add(mgr.cfg.TimeOut)
	s.Cookie.Expires = s.Expires
	http.SetCookie(s._w, s.Cookie)
}

func generateUUID() string {
	return fmt.Sprintf("%s-%s", uuid.New().String(), uuid.New().String())
}

func (s *Session) copy(cookie *http.Cookie) error {
	s.Cookie = new(http.Cookie)
	c, err := deepcopy.Anything(cookie)
	if err != nil {
		return errors.New("cookie make a deep copy from src into dst failed")
	}
	s.Cookie = c.(*http.Cookie)
	return nil
}

func (s *Session) MigrateSession() error {

	s.ID = generateUUID()
	newSession, err := deepcopy.Anything(s)
	if err != nil {
		return errors.New("migrate session make a deep copy from src into dst failed")
	}

	newSession.(*Session).ID = s.ID
	newSession.(*Session).Cookie.Value = s.ID
	newSession.(*Session).Expires = time.Now().Add(mgr.cfg.TimeOut)
	newSession.(*Session)._w = s._w
	newSession.(*Session).refreshCookie()

	return mgr.store.Create(newSession.(*Session))
}

// It makes a deep copy by using json.Marshal and json.Unmarshal, so it's not very
// performant.
// Make a deep copy from src into dst.

// fix bug:
// 1.Types of function parameters can be combined
// https://deepsource.io/gh/higker/sessionx/issue/CRT-A0017/occurrences
// 2.Incorrectly formatted error string
// https://deepsource.io/gh/higker/sessionx/issue/SCC-ST1005/occurrences
func _copy(dst, src interface{}) error {
	if dst == nil {
		return fmt.Errorf("dst cannot be nil")
	}
	if src == nil {
		return fmt.Errorf("src cannot be nil")
	}
	bytes, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("unable to marshal src: %s", err)
	}
	err = json.Unmarshal(bytes, dst)
	if err != nil {
		return fmt.Errorf("unable to unmarshal into dst: %s", err)
	}
	return nil
}
