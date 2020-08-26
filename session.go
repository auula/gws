// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 9:10 PM - UTC/GMT+08:00

package session

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Session interface {
	Get(key string) ([]byte, error)
	Set(key string, data interface{}) error
	Del(key string) error
	Clean(w http.ResponseWriter)
}

// Session for  Memory item
type MemorySession struct {
	ID      string                 // Unique id
	Safe    sync.Mutex             // Mutex lock
	Expires time.Time              // Expires time
	Data    map[string]interface{} // Save data
}

//实例化
func newMSessionItem(id string, maxAge int) *MemorySession {
	return &MemorySession{
		Data:    make(map[string]interface{}, maxSize),
		ID:      id,
		Expires: time.Now().Add(time.Duration(maxAge) * time.Second),
	}
}

// Builder build  session store
func Builder(store StoreType, conf *Config) error {
	if conf.MaxAge < DefaultMaxAge {
		return errors.New("session maxAge no less than 30min")
	}

	_Cfg = conf
	switch store {
	default:
		return errors.New("build session error, not implement type store")
	case Memory:
		_Store = newMemoryStore()
		_Cfg._st = Memory
		return nil
	case Redis:
		//redisStore, err := newRedisStore()
		//if err != nil {
		//	return err
		//}
		//_Store = redisStore
		_Cfg._st = Redis
		return nil
	}
}

// Ctx return request session object
func Ctx(writer http.ResponseWriter, request *http.Request) (Session, error) {

	// 检测是否有这个session数据
	// 1.全局session垃圾回收器
	// 2.请求一过来就进行检测
	// 3.如果请求的id值在内存存在并且垃圾回收也存在时间也是有效的就说明这个session是有效的

	cookie, err := request.Cookie(_Cfg.CookieName)
	if _Cfg._st == Memory {
		if err != nil || len(cookie.Value) <= 0 {
			item := recordCookie(writer, cookie)
			return item, nil
		}
		// 防止浏览器关闭重新打开抛异常
		sid, err := url.QueryUnescape(cookie.Value)
		if err != nil {
			return nil, err
		}
		session := _Store.(*MemoryStore).values[sid]
		if session == nil {
			item := recordCookie(writer, cookie)
			_Store.(*MemoryStore).values[sid] = item.(*MemorySession)
			return item, nil
		}
		return _Store.(*MemoryStore).values[cookie.Value], nil
	}

	return nil, nil
}

// Get get session data by key
func (ms *MemorySession) Get(key string) ([]byte, error) {
	if key == "" || len(key) <= 0 {
		return nil, ErrorKeyNotExist
	}
	//var result Value
	//result.Key = key
	// 把id和到期时间传过去方便后面使用
	cv := map[string]interface{}{contextValueID: ms.ID, contextValueKey: key}
	value := context.WithValue(context.TODO(), contextValue, cv)
	b, err := _Store.Reader(value)
	if err != nil {
		return nil, err
	}
	//result.Value = b
	return b, nil
}

// Set set session data by key
func (ms *MemorySession) Set(key string, data interface{}) error {
	if key == "" || len(key) <= 0 {
		return ErrorKeyFormat
	}
	cv := map[string]interface{}{contextValueID: ms.ID, contextValueKey: key, contextValueData: data}
	value := context.WithValue(context.TODO(), contextValue, cv)
	return _Store.Writer(value)
}

// Del delete session data by key
func (ms *MemorySession) Del(key string) error {
	if key == "" || len(key) <= 0 {
		return ErrorKeyFormat
	}
	cv := map[string]interface{}{contextValueID: ms.ID, contextValueKey: key}
	value := context.WithValue(context.TODO(), contextValue, cv)
	_Store.Remove(value)
	return nil
}

// Clean clean session data
func (ms *MemorySession) Clean(w http.ResponseWriter) {
	cv := map[string]interface{}{contextValueID: ms.ID}
	value := context.WithValue(context.TODO(), contextValue, cv)
	_Store.Clean(value)
	cookie := &http.Cookie{
		Name:     _Cfg.CookieName,
		Value:    "",
		Path:     _Cfg.Path,
		Domain:   _Cfg.Domain,
		Secure:   _Cfg.Secure,
		MaxAge:   -1,
		Expires:  time.Now().AddDate(-1, 0, 0),
		HttpOnly: _Cfg.HttpOnly,
	}
	http.SetCookie(w, cookie)
}

// 检测sessionID是否有效
func checkID(id string) bool {
	if len(_Store.(*MemoryStore).values) > 0 {
		return _Store.(*MemoryStore).values[id] == nil
	}
	return false
}

func recordCookie(w http.ResponseWriter, cookie *http.Cookie) Session {
	// 创建一个cookie
	sid := string(Random(32, RuleKindAll))
	cookie = &http.Cookie{
		Name: _Cfg.CookieName,
		//这里是并发不安全的，但是这个方法已上锁
		Value:    url.QueryEscape(sid), //转义特殊符号@#￥%+*-等
		Path:     _Cfg.Path,
		Domain:   _Cfg.Domain,
		HttpOnly: _Cfg.HttpOnly,
		Secure:   _Cfg.Secure,
		MaxAge:   int(_Cfg.MaxAge),
		Expires:  time.Now().Add(time.Duration(_Cfg.MaxAge)),
	}
	http.SetCookie(w, cookie) //设置到响应中
	item := newMSessionItem(sid, int(_Cfg.MaxAge))
	_Store.(*MemoryStore).values[sid] = item
	return item
}
