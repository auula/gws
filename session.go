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
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"net/http"
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
	_w      http.ResponseWriter
}

// Get Retrieves the stored element data from the session via the key
func (s *Session) Get(key string) (interface{}, error) {
	err := mgr.store.Reader(s)
	s.refreshCookie()
	if err != nil {
		return nil, err
	}
	if ele, ok := s.Data[key]; ok {
		return ele, nil
	}
	return nil, fmt.Errorf("key '%s' does not exist", key)
}

// Set Stores information in the session
func (s *Session) Set(key string, v interface{}) error {
	mux.Lock()
	if s.Data == nil {
		s.Data = make(map[string]interface{}, 8)
	}
	s.Data[key] = v
	mux.Unlock()
	s.refreshCookie()
	return mgr.store.Update(s)
}

// Remove an element stored in the session
func (s *Session) Remove(key string) error {
	s.refreshCookie()
	return mgr.store.Remove(s, key)
}

// Clean up all data for this session
func (s *Session) Clean() error {
	return mgr.store.Delete(s)
}

// Handler Get session data from the Request
func Handler(w http.ResponseWriter, req *http.Request) *Session {
	mux.Lock()
	defer mux.Unlock()

	// 从请求里面取session
	var session Session
	session._w = w
	cookie, err := req.Cookie(mgr.cfg.Cookie.Name)
	if err != nil || cookie == nil || len(cookie.Value) <= 0 {
		return createSession(w, cookie, &session)
	}

	// ID通过编码之后长度是48位
	if len(cookie.Value) >= 48 {
		sDec, _ := base64.StdEncoding.DecodeString(cookie.Value)
		session.ID = string(sDec)
		if mgr.store.Reader(&session) != nil {
			return createSession(w, cookie, &session)
		}
	}
	return &session
}

func createSession(w http.ResponseWriter, cookie *http.Cookie, session *Session) *Session {
	sessionId := uuid.New().String()
	expireTime := time.Now().Add(mgr.cfg.TimeOut)

	// init cookie parameter
	cookie = mgr.cfg.Cookie
	cookie.Expires = expireTime
	cookie.Value = base64.StdEncoding.EncodeToString([]byte(sessionId))

	// init session parameter
	session.ID = sessionId
	session.Expires = expireTime
	_ = mgr.store.Create(session)
	http.SetCookie(w, cookie)
	return session
}

// 刷新cookie 会话只要有操作就重置会话生命周期
func (s *Session) refreshCookie() {
	mgr.cfg.Cookie.Expires = time.Now().Add(mgr.cfg.TimeOut)
	http.SetCookie(s._w, mgr.cfg.Cookie)
}
