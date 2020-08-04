// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/4 - 10:40 PM

package session

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

type StoreType int8

const (
	MemoryStore StoreType = iota
	RedisStore
	maxSize = 16
)

var (
	_Store            Store
	DefaultCookieName = "go_session_key"
	DefaultMaxAge     = 60 * 30 // 30 min
)

// Session standard
type Session interface {
	Set(Key, value interface{})
	Get(Key interface{}) interface{}
	Remove(Key interface{})
	ID() string
	Clear()
}

// Config param
type Config struct {
	// cookie参数
	CookieName string // sessionID的cookie键名
	Domain     string // sessionID的cookie作用域名
	Path       string // sessionID的cookie作用路径
	//Key                       string        // sessionID值加密的密钥
	RedisAddr                 string        // redis地址
	RedisPassword             string        // redis密码
	RedisKeyPrefix            string        // redis键名前缀
	MaxAge                    int           // 最大生命周期（秒）
	IdleTime                  time.Duration // 空闲生命周期
	RedisDB                   int           // redis数据库
	HttpOnly                  bool          // 仅用于http（无法被js读取）
	Secure                    bool          // 启用https
	DisableAutoUpdateIdleTime bool          // 禁止自动更新空闲时间
}

// Item a session  item
type Item struct {
	SID              string                      // unique id
	Safe             sync.Mutex                  // mutex lock
	LastAccessedTime time.Time                   // last visit time
	MaxAge           int64                       // over time
	Data             map[interface{}]interface{} // save data
}

//实例化
func newSessionItem(id string) *Item {
	return &Item{
		Data: make(map[interface{}]interface{}, maxSize),
		SID:  id,
	}
}

//同一个会话均可调用，进行设置，改操作必须拥有排斥锁
func (si *Item) Set(key, value interface{}) {
	si.Safe.Lock()
	defer si.Safe.Unlock()
	si.Data[key] = value
}

func (si *Item) Get(key interface{}) interface{} {
	if value := si.Data[key]; value != nil {
		return value
	}
	return nil
}

func (si *Item) Remove(key interface{}) {
	if value := si.Data[key]; value != nil {
		delete(si.Data, key)
	}
}

func (si *Item) ID() string {
	return si.SID
}
func (si *Item) Clear() {

}

func Builder(store StoreType, conf *Config) (Session, error) {
	switch store {
	default:
		return nil, errors.New("build session error, not implement type store")
	case MemoryStore:

	case RedisStore:

	}
	return nil, nil
}

// Action Session handle func
func Action(w http.ResponseWriter, r *http.Request) (Session, error) {

	return nil, nil
}
