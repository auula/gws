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
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	mux sync.Mutex
	mgr *manager
)

type Session struct {
	// 会话ID
	ID string
	// session超时时间
	Expires time.Time
	// 存储数据的map
	Data map[string]interface{}
	_w   http.ResponseWriter
	// 每个session对应一个cookie
	cookie *http.Cookie
}

// Get Retrieves the stored element data from the session via the key
func (s *Session) Get(key string) (interface{}, error) {
	err := mgr.store.Reader(s)
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
	s.refreshCookie()
	return mgr.store.Delete(s)
}

// Handler Get session data from the Request
func Handler(w http.ResponseWriter, req *http.Request) *Session {
	mux.Lock()
	defer mux.Unlock()
	// 从请求里面取session
	var session Session
	cookie, err := req.Cookie(mgr.cfg.Cookie.Name)
	if err != nil || cookie == nil || len(cookie.Value) <= 0 {
		return createSession(w, cookie, &session)
	}
	// ID通过编码之后长度是73位
	if len(cookie.Value) >= 73 {
		session.ID = cookie.Value
		if mgr.store.Reader(&session) != nil {
			return createSession(w, cookie, &session)
		}
		// 防止web服务器重启之后redis会话数据还在
		// 但是浏览器cookie没有更新
		// 重新刷新cookie
		mgr.cfg.Cookie.Value = session.ID
		mgr.cfg.Cookie.Expires = session.Expires
		session.cookie = mgr.cfg.Cookie
		http.SetCookie(w, session.cookie)
	}
	return &session
}

func createSession(w http.ResponseWriter, cookie *http.Cookie, session *Session) *Session {
	sessionId := generateUUID()
	// init session parameter
	session.ID = sessionId
	session.Expires = time.Now().Add(mgr.cfg.TimeOut)
	_ = mgr.store.Create(session)

	// init cookie parameter
	session.cookie = mgr.cfg.Cookie
	session.cookie.Expires = session.Expires
	session.cookie.Value = sessionId

	http.SetCookie(w, session.cookie)
	return session
}

// 刷新cookie 会话只要有操作就重置会话生命周期
func (s *Session) refreshCookie() {
	s.Expires = time.Now().Add(mgr.cfg.TimeOut)
	s.cookie.Expires = s.Expires
	// 这里不是使用指针
	// 因为这里我们支持redis 如果web服务器重启了
	// 那么session数据在内存里清空
	// 从redis读取的数据反序列化地址和重新启动的不一样
	// 所有直接数据拷贝
	http.SetCookie(s._w, s.cookie)
}

func generateUUID() string {
	return fmt.Sprintf("%s-%s", uuid.New().String(), uuid.New().String())
}
