// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/7/25 - 6:08 PM

package session

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Manager session manager
type Manager struct {
	// cookie name
	CookieName string
	// data storage
	Store Storage
	// max life cycle
	MaxAge int64
	// safe lock
	Safe sync.Mutex
}

//实例化一个session管理器
func New(storeType StoreType, cookieName string, maxAge int64) *Manager {

	var _store *MemoryStore
	switch storeType {
	case MemoryType:
		_store = newMemoryStore()
	case RedisType:
		panic("not implement store type!")
	case DatabaseType:
		panic("not implement store type!")
	default:
		panic("not implement store type!")
	}

	sessionManager := &Manager{
		CookieName: cookieName,
		Store:      _store, //默认以内存实现
		MaxAge:     maxAge, //默认30分钟
	}

	go sessionManager.GC()

	return sessionManager
}

func (m *Manager) GetCookieName() string {
	return m.CookieName
}

//先判断当前请求的cookie中是否存在有效的session,存在返回，不存在创建
func (m *Manager) BeginSession(w http.ResponseWriter, r *http.Request) Session {
	//防止处理时，进入另外的请求
	m.Safe.Lock()
	defer m.Safe.Unlock()
	cookie, err := r.Cookie(m.CookieName)
	if err != nil || cookie.Value == "" { //如果当前请求没有改cookie名字对应的cookie
		//创建一个
		sid := m.randomId()
		//根据保存session方式，如内存，数据库中创建
		session, _ := m.Store.InitSession(sid, m.MaxAge) //该方法有自己的锁，多处调用到

		maxAge := m.MaxAge

		if maxAge == 0 {
			maxAge = session.(*Memory).MaxAge
		}
		//用session的ID于cookie关联
		//cookie名字和失效时间由session管理器维护
		cookie := http.Cookie{
			Name: m.CookieName,
			//这里是并发不安全的，但是这个方法已上锁
			Value:    url.QueryEscape(sid), //转义特殊符号@#￥%+*-等
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(maxAge),
			Expires:  time.Now().Add(time.Duration(maxAge)),
		}
		http.SetCookie(w, &cookie) //设置到响应中
		return session
	} else { //如果存在
		sid, _ := url.QueryUnescape(cookie.Value)       //反转义特殊符号
		session := m.Store.(*MemoryStore).sessions[sid] //从保存session介质中获取
		if session == nil {
			//创建一个
			//sid := m.randomId()
			//根据保存session方式，如内存，数据库中创建
			newSession, _ := m.Store.InitSession(sid, m.MaxAge) //该方法有自己的锁，多处调用到

			maxAge := m.MaxAge

			if maxAge == 0 {
				maxAge = newSession.(*Memory).MaxAge
			}
			//用session的ID于cookie关联
			//cookie名字和失效时间由session管理器维护
			newCookie := http.Cookie{
				Name: m.CookieName,
				//这里是并发不安全的，但是这个方法已上锁
				Value:    url.QueryEscape(sid), //转义特殊符号@#￥%+*-等
				Path:     "/",
				HttpOnly: true,
				MaxAge:   int(maxAge),
				Expires:  time.Now().Add(time.Duration(maxAge)),
			}
			http.SetCookie(w, &newCookie) //设置到响应中
			return newSession
		}
		return session
	}

}

//更新超时
func (m *Manager) Update(w http.ResponseWriter, r *http.Request) {
	m.Safe.Lock()
	defer m.Safe.Unlock()

	cookie, err := r.Cookie(m.CookieName)
	if err != nil {
		return
	}
	t := time.Now()
	sid, _ := url.QueryUnescape(cookie.Value)

	sessions := m.Store.(*MemoryStore).sessions
	session := sessions[sid].(*Memory)
	session.LastAccessedTime = t
	sessions[sid] = session

	if m.MaxAge != 0 {
		cookie.MaxAge = int(m.MaxAge)
	} else {
		cookie.MaxAge = int(session.MaxAge)
	}
	http.SetCookie(w, cookie)
}

//通过ID获取session
func (m *Manager) GetSessionById(sid string) Session {
	session := m.Store.(*MemoryStore).sessions[sid]
	return session
}

//是否内存中存在
func (m *Manager) MemoryIsExists(sid string) bool {
	_, ok := m.Store.(*MemoryStore).sessions[sid]
	return ok
}

//手动销毁session，同时删除cookie
func (m *Manager) Destroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(m.CookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		m.Safe.Lock()
		defer m.Safe.Unlock()

		sid, _ := url.QueryUnescape(cookie.Value)
		_ = m.Store.DestroySession(sid)

		cookie2 := http.Cookie{
			MaxAge:  0,
			Name:    m.CookieName,
			Value:   "",
			Path:    "/",
			Expires: time.Now().Add(time.Duration(0)),
		}

		http.SetCookie(w, &cookie2)
	}
}

func (m *Manager) CookieIsExists(r *http.Request) bool {
	_, err := r.Cookie(m.CookieName)
	if err != nil {
		return false
	}
	return true
}

//开启每个会话，同时定时调用该方法
//到达session最大生命时，且超时时。回收它
func (m *Manager) GC() {
	m.Safe.Lock()
	defer m.Safe.Unlock()

	m.Store.GCSession()
	//在多长时间后执行匿名函数，这里指在某个时间后执行GC
	time.AfterFunc(time.Duration(m.MaxAge*10), func() {
		m.GC()
	})
}

//是否将session放入内存（操作内存）默认是操作内存
func (m *Manager) IsFromMemory() {
	m.Store = newMemoryStore()
}

//是否将session放入数据库（操作数据库）
func (m *Manager) IsFromDB() {
	//TODO
	//关于存数据库暂未实现
}

func (m *Manager) SetMaxAge(t int64) {
	m.MaxAge = t
}

//如果你自己实现保存session的方式，可以调该函数进行定义
func (m *Manager) SetSessionFrom(storage Storage) {
	m.Store = storage
}

//生成一定长度的随机数
func (m *Manager) randomId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	//加密
	return base64.URLEncoding.EncodeToString(b)
}
