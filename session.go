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

// Ctx return request session object
func Ctx(writer http.ResponseWriter, request *http.Request) (*Session, error) {
	var session *Session
	cookie, err := request.Cookie(_Cfg.CookieName)
	if err != nil || cookie == nil || len(cookie.Value) <= 0 {
		// 重新生成一个cookie 和唯一 sessionID
		nc := newCookie(writer)
		unescape, err := url.QueryUnescape(nc.Value)
		if err != nil {
			return nil, err
		}
		// 把加密的id解析成存储sid
		sid, err := decodeStorageID(unescape, _Cfg.EncryptedKey)
		if err != nil {
			return nil, err
		}
		session = &Session{
			Cookie: nc,
			ID:     sid,
		}
		return session, nil
	}
	// 防止浏览器关闭重新打开抛异常
	id, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return nil, err
	}
	// 把加密的id解析成存储sid
	sid, err := decodeStorageID(id, _Cfg.EncryptedKey)
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
	// 把加密的id解析成存储sid
	sid, err := decodeStorageID(s.ID, _Cfg.EncryptedKey)
	if err != nil {
		return nil, err
	}
	b, err := _Store.Reader(sid, key)
	if err != nil {
		return nil, err
	}
	//result.Value = b
	return b, nil
}

// Set set session data by key
func (s *Session) Set(key string, data interface{}) error {
	if key == "" || len(key) <= 0 {
		return ErrorKeyFormat
	}
	// 把加密的id解析成存储sid
	sid, err := decodeStorageID(s.ID, _Cfg.EncryptedKey)
	if err != nil {
		return err
	}
	return _Store.Writer(sid, key, data)
}

// Del delete session data by key
func (s *Session) Del(key string) error {
	if key == "" || len(key) <= 0 {
		return ErrorKeyFormat
	}
	// 把加密的id解析成存储sid
	sid, err := decodeStorageID(s.ID, _Cfg.EncryptedKey)
	if err != nil {
		return err
	}
	_Store.Remove(sid, key)
	return nil
}

// Clean clean session data
func (s *Session) Clean() {
	// 把加密的id解析成存储sid
	sid, _ := decodeStorageID(s.ID, _Cfg.EncryptedKey)
	_Store.Clean(sid)
}

func newCookie(w http.ResponseWriter) *http.Cookie {
	// 创建一个cookie
	s := Random(32, RuleKindAll)
	// 把存储id加密返回
	bytes, _ := encodeByBytes(strToByte(_Cfg.EncryptedKey), s)
	cookie := &http.Cookie{
		Name: _Cfg.CookieName,
		//这里是并发不安全的，但是这个方法已上锁
		Value:    url.QueryEscape(bytes), //转义特殊符号@#￥%+*-等
		Path:     _Cfg.Path,
		Domain:   _Cfg.Domain,
		HttpOnly: _Cfg.HttpOnly,
		Secure:   _Cfg.Secure,
		MaxAge:   int(_Cfg.MaxAge),
		Expires:  time.Now().Add(time.Duration(_Cfg.MaxAge)),
	}
	http.SetCookie(w, cookie) //设置到响应中
	fmt.Println(cookie)
	return cookie
}
