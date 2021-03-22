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
	"net/http"
	"runtime"
	"time"

	"github.com/google/uuid"
)

type storeType uint8

const (
	// memoryStore store type
	M storeType = iota
	// redis store type
	R
	SESSION_KEY = "sessionx-id"
)

// manager for session manager
type manager struct {
	cfg   *Configs
	store storage
}

func New(t storeType, cfg *Configs) {
	switch t {
	case M:
		m := new(memoryStore)
		m.sessions = make(map[string]*Session, 512*runtime.NumCPU())
		go m.GC()
		mgr = &manager{cfg: cfg, store: m}
	case R:
		panic("not impl store type")
	default:
		panic("not impl store type")
	}
}

func Handler(w http.ResponseWriter, req *http.Request) *Session {
	mux.Lock()
	defer mux.Unlock()
	var session Session
	cookie, err := req.Cookie(mgr.cfg.Cookie.Name)
	if err != nil || cookie == nil || len(cookie.Value) <= 0 {
		return createSession(w, cookie, &session)
	}
	if len(cookie.Value) >= 48 {
		sDec, _ := base64.StdEncoding.DecodeString(cookie.Value)
		session.ID = string(sDec)
		reader, err := mgr.store.Reader(&session)
		if err != nil {
			return createSession(w, cookie, &session)
		}
		_ = decoder(reader, &session)
	}
	return &session
}

func createSession(w http.ResponseWriter, cookie *http.Cookie, session *Session) *Session {
	sessionId := uuid.New().String()
	expireTime := time.Now().Add(mgr.cfg.SessionLifeTime)

	// init cookie parameter
	cookie = mgr.cfg.Cookie
	cookie.Expires = expireTime
	cookie.Value = base64.StdEncoding.EncodeToString([]byte(sessionId))

	// init session parameter
	session.ID = sessionId
	session.Expires = expireTime
	_, _ = mgr.store.Create(session)
	http.SetCookie(w, cookie)
	return session
}
