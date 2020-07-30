// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/7/24 - 07:47 PM

package session

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

var (
	_Store     Store
	CookieName = "go_session_key"
	MaxAge     = 60 * 30 // 30 min
)

// Session Operation interface, Session operation of different storage methods is different,
// and the implementation is also different.
type Session interface {
	Set(key, value interface{})
	Get(key interface{}) interface{}
	Remove(key interface{}) error
	GetID() string
}

// Builder build session
func Builder(storeType Store) error {
	switch storeType.(type) {
	default:
		return errors.New("not implement store type")
	case *MemoryStore:
		_Store = storeType
		go storeType.(*MemoryStore).GC()
	}
	return nil
}

func Handle(w http.ResponseWriter, r *http.Request) (Session, error) {
	//防止处理时，进入另外的请求
	cookie, err := r.Cookie(CookieName)
	if err != nil || len(cookie.Value) <= 0 {
		_, item := setCookie(w, cookie)
		return item, nil
	}
	// 防止浏览器关闭重新打开抛异常
	sid, _ := url.QueryUnescape(cookie.Value)
	session := _Store.(*MemoryStore).sessions[sid]
	if session == nil {
		_, item := setCookie(w, cookie)
		return item, nil
	}
	return _Store.(*MemoryStore).sessions[cookie.Value], nil
}

func setCookie(w http.ResponseWriter, cookie *http.Cookie) (*http.Cookie, *MemoryItem) {
	// 创建一个cookie
	sid := string(Random(32, RuleKindAll))
	cookie = &http.Cookie{
		Name: CookieName,
		//这里是并发不安全的，但是这个方法已上锁
		Value:    url.QueryEscape(sid), //转义特殊符号@#￥%+*-等
		Path:     "/",
		HttpOnly: true,
		MaxAge:   MaxAge,
		Expires:  time.Now().Add(time.Duration(MaxAge)),
	}
	http.SetCookie(w, cookie) //设置到响应中
	item := newMemoryItem(sid)
	_Store.(*MemoryStore).sessions[sid] = item
	return cookie, item
}
