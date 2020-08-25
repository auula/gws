// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 9:10 PM - UTC/GMT+08:00

package session

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Session Unite struct
type Session struct {
	ID     string
	Cookie *http.Cookie
}

// Builder build  session store
func Builder(store StoreType, conf *Config) error {
	//if conf.MaxAge < DefaultMaxAge {
	//	return errors.New("session maxAge no less than 30min")
	//}
	_Cfg = conf
	fmt.Println(conf)
	switch store {
	default:
		return errors.New("build session error, not implement type store")
	case Memory:
		_Store = newMemoryStore()
		return nil
	case Redis:
		redisStore, err := newRedisStore()
		if err != nil {
			return err
		}
		_Store = redisStore
		return nil
	}
}

// Capture return request session object
func Capture(writer http.ResponseWriter, request *http.Request) (*Session, error) {
	var session *Session
	cookie, err := request.Cookie(_Cfg.CookieName)
	if err != nil || cookie == nil || len(cookie.Value) <= 0 {
		session = &Session{Cookie: cookie, ID: string(Random(16, RuleKindAll))}
		recordCookie(writer, session)
		return session, nil
	}

	// 防止浏览器关闭重新打开抛异常
	sid, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return nil, err
	}
	session = &Session{ID: sid, Cookie: cookie}
	return session, nil
}

// Get get session data by key
func (s *Session) Get(key string) ([]byte, error) {
	if key == "" || len(key) <= 0 {
		return nil, ErrorKeyNotExist
	}
	//var result Value
	//result.Key = key
	b, err := _Store.Reader(s.ID, key)
	if err != nil {
		return nil, err
	}
	fmt.Println(b)
	fmt.Println(string(b))
	//result.Value = b
	return b, nil
}

// Set set session data by key
func (s *Session) Set(key string, data interface{}) error {
	if key == "" || len(key) <= 0 {
		return ErrorKeyFormat
	}
	return _Store.Writer(s.ID, key, data)
}

// Del delete session data by key
func (s *Session) Del(key string) error {
	if key == "" || len(key) <= 0 {
		return ErrorKeyFormat
	}
	_Store.Remove(s.ID, key)
	return nil
}

// Clean clean session data
func (s *Session) Clean() {
	_Store.Clean(s.ID)
}

func recordCookie(w http.ResponseWriter, session *Session) {
	// 创建一个cookie
	cookie := &http.Cookie{
		Name: _Cfg.CookieName,
		//这里是并发不安全的，但是这个方法已上锁
		Value:    url.QueryEscape(session.ID), //转义特殊符号@#￥%+*-等
		Path:     _Cfg.Path,
		Domain:   _Cfg.Domain,
		HttpOnly: _Cfg.HttpOnly,
		Secure:   _Cfg.Secure,
		MaxAge:   int(_Cfg.MaxAge),
		Expires:  time.Now().Add(time.Duration(_Cfg.MaxAge)),
	}
	http.SetCookie(w, cookie) //设置到响应中
}
