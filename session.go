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
	"github.com/barkimedes/go-deepcopy"
	"github.com/google/uuid"
	"net/http"
	"sync"
	"time"
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
	// 会话ID
	ID string
	// session超时时间
	Expires time.Time
	// 存储数据的map
	Data map[interface{}]interface{}
	_w   http.ResponseWriter
	// 每个session对应一个cookie
	Cookie *http.Cookie
}

// Get Retrieves the stored element data from the session via the key
func (s *Session) Get(key interface{}) (interface{}, error) {
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
func (s *Session) Set(key, v interface{}) error {

	lock["W"](func() {
		if s.Data == nil {
			s.Data = make(map[interface{}]interface{}, 8)
		}
		s.Data[key] = v
	})

	s.refreshCookie()
	return mgr.store.Update(s)
}

// Remove an element stored in the session
func (s *Session) Remove(key interface{}) error {
	s.refreshCookie()

	lock["R"](func() {
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
	// 从请求里面取session
	var session Session
	session._w = w
	cookie, err := req.Cookie(mgr.cfg.Cookie.Name)
	if err != nil || cookie == nil || len(cookie.Value) <= 0 {
		return createSession(w, cookie, &session)
	}
	// ID通过编码之后长度是73位
	if len(cookie.Value) >= 73 {
		session.ID = cookie.Value
		if mgr.store.Read(&session) != nil {
			return createSession(w, cookie, &session)
		}

		// 防止web服务器重启之后redis会话数据还在
		// 但是浏览器cookie没有更新
		// 重新刷新cookie

		// 存在指针一致问题，这样操作还是一块内存，所有我们需要复制副本
		_ = session.copy(mgr.cfg.Cookie)
		session.Cookie.Value = session.ID
		session.Cookie.Expires = session.Expires
		http.SetCookie(w, session.Cookie)
	}
	// 地址一样不行！！！
	// log.Printf("mgr.cfg.Cookie pointer:%p \n", mgr.cfg.Cookie)
	// log.Printf("session.cookie pointer:%p \n", session.Cookie)
	return &session
}

func createSession(w http.ResponseWriter, cookie *http.Cookie, session *Session) *Session {
	// init session parameter
	session.ID = generateUUID()
	session.Expires = time.Now().Add(mgr.cfg.TimeOut)
	_ = mgr.store.Create(session)

	// 重置配置cookie模板
	session.copy(mgr.cfg.Cookie)
	session.Cookie.Value = session.ID
	session.Cookie.Expires = session.Expires

	http.SetCookie(w, session.Cookie)
	return session
}

// 刷新cookie 会话只要有操作就重置会话生命周期
func (s *Session) refreshCookie() {
	s.Expires = time.Now().Add(mgr.cfg.TimeOut)
	s.Cookie.Expires = s.Expires
	// 这里不是使用指针
	// 因为这里我们支持redis 如果web服务器重启了
	// 那么session数据在内存里清空
	// 从redis读取的数据反序列化地址和重新启动的不一样
	// 所有直接数据拷贝
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
	// 迁移到新内存 防止会话一致引发安全问题
	// 这个问题的根源在 sessionid 不变，如果用户在未登录时拿到的是一个 sessionid，登录之后服务端给用户重新换一个 sessionid，就可以防止会话固定攻击了。
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
	// 新内存开始持久化
	// log.Printf("old session pointer:%p \n", s)
	// log.Printf("new session pointer:%p \n", newSession.(*Session))
	//log.Println("MigrateSession:", newSession.(*Session))
	return mgr.store.Create(newSession.(*Session))
}

// 为什么选择使用Data作为k v存储 而不是像持久化存储一样用sync.map

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
