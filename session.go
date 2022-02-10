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

	// Global Configure controller
	globalConfig *Config

	// Session concurrent safe mutex
	migrateMux sync.Mutex

	ErrKeyNoData          = errors.New("key no data")
	ErrSessionNoData      = errors.New("session no data")
	ErrIsEmpty            = errors.New("key or session id is empty")
	ErrAlreadyExpired     = errors.New("session already expired")
	ErrRemoveSessionFail  = errors.New("remove session fail")
	ErrMigrateSessionFail = errors.New("migrate session fail")
)

// Values is session item value
type Values map[string]interface{}

// Session is web session struct
type Session struct {
	session
}

// session struct
type session struct {
	id         string
	CreateTime time.Time
	ExpireTime time.Time
	Values
}

// GetSession Get session data from the Request
func GetSession(w http.ResponseWriter, req *http.Request) (*Session, error) {
	var session Session

	cookie, err := req.Cookie(globalConfig.CookieName)
	if cookie == nil || err != nil {
		debug.trace(cookie)
		return createSession(w, cookie)
	}

	if len(cookie.Value) >= 73 {
		session.id = cookie.Value
		if globalStore.Read(&session) != nil {
			return createSession(w, cookie)
		}
	}

	debug.trace(&session)
	return &session, nil
}

// ID return session id
func (s *Session) ID() string {
	return s.id
}

// Sync save data modify
func (s *Session) Sync() error {
	debug.trace(s)
	return globalStore.Write(s)
}

// Migrate migrate old session data to new session
func Migrate(write http.ResponseWriter, old *Session) (*Session, error) {
	var (
		ns     = NewSession()
		cookie = NewCookie()
	)

	migrateMux.Lock()
	ns.Values = old.Values
	cookie.Value = ns.id
	cookie.MaxAge = int(globalConfig.LifeTime) / 1e9
	migrateMux.Unlock()

	return ns,
		func() error {
			if ns.Sync() != nil {
				return ErrMigrateSessionFail
			}
			if globalStore.Remove(old) != nil {
				return ErrRemoveSessionFail
			}
			http.SetCookie(write, cookie)
			return nil
		}()
}

// createSession return new session
func createSession(w http.ResponseWriter, cookie *http.Cookie) (*Session, error) {

	// FIX BUG:
	// https://deepsource.io/gh/auula/gws/run/5b13c99b-9101-4e4f-8197-acfd730c28a0/go/SCC-SA4009
	session := NewSession()

	debug.trace(session)

	if cookie == nil {
		cookie = NewCookie()
	}
	cookie.Value = session.id
	cookie.MaxAge = int(globalConfig.LifeTime) / 1e9
	if err := globalStore.Write(session); err != nil {
		return nil, err
	}

	debug.trace(cookie)

	http.SetCookie(w, cookie)

	debug.trace(session)
	return session, nil
}

// NewCookie return default config cookie pointer
func NewCookie() *http.Cookie {
	return &http.Cookie{
		Domain:   globalConfig.Domain,
		Path:     globalConfig.Path,
		Name:     globalConfig.CookieName,
		Secure:   globalConfig.Secure,
		HttpOnly: globalConfig.HttpOnly,
	}
}

// uuid73 generate session uuid length 73
func uuid73() string {
	return fmt.Sprintf("%s-%s", uuid.New().String(), uuid.New().String())
}

// NewSession return new session
func NewSession() *Session {
	nowTime := time.Now()
	return &Session{
		session: session{
			id:         uuid73(),
			Values:     make(Values),
			CreateTime: nowTime,
			ExpireTime: nowTime.Add(lifeTime),
		},
	}
}

// Expired check current session whether expire
func (s *Session) Expired() bool {
	return time.Duration(s.ExpireTime.UnixNano()) <= time.Duration(time.Now().UnixNano())
}

// Invalidate remove the session
func Invalidate(s *Session) error {
	debug.trace(s)
	return globalStore.Remove(s)
}

// Malloc reallocation of memory
func Malloc(v *Values) {
	*v = make(Values)
}

// Open Initialize storage with custom configuration
func Open(opt Configure) {
	debug.trace(opt)

	globalConfig = opt.Parse()

	switch globalConfig.store {
	case ram:
		globalStore = NewRAM()
	case rds:
		rdb := NewRds()
		timeout, cancelFunc := timeoutCtx()
		defer cancelFunc()
		if err := rdb.store.Ping(timeout).Err(); err != nil {
			panic(err.Error())
		}
		globalStore = rdb
	default:
		globalStore = NewRAM()
	}
}

// StoreFactory Initialize custom storage media
func StoreFactory(opt Options, store Storage) {
	globalConfig = opt.Parse()
	globalStore = store
}
