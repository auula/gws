// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/4 - 10:40 PM

package session

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

type StoreType int8

const (
	MemoryStore StoreType = iota
	RedisStore
	maxSize           = 16
	DefaultCookieName = "go_session_key"
	DefaultMaxAge     = 30 * 60 // 30 min
)

var (
	_Store Store
	_Cfg   *Config
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
	MaxAge     int    // 最大生命周期（秒）
	HttpOnly   bool   // 仅用于http（无法被js读取）
	Secure     bool   // 启用https
	//Key                       string        // sessionID值加密的密钥
	//RedisAddr      string // redis地址
	//RedisPassword  string // redis密码
	//RedisKeyPrefix string // redis键名前缀
	// IdleTime                  time.Duration // 空闲生命周期
	// RedisDB                   int           // redis数据库
	// DisableAutoUpdateIdleTime bool          // 禁止自动更新空闲时间
}

// Builder build  session store
func Builder(store StoreType, conf *Config) error {
	if conf.MaxAge < DefaultMaxAge {
		return errors.New("session maxAge no less than 30min")
	}
	switch store {
	default:
		return errors.New("build session error, not implement type store")
	case MemoryStore:
		_Cfg = conf
		_Store = newMemoryStore()
		go _Store.GC()
		return nil
	case RedisStore:
		return errors.New("not implement type store,to github: https://github.com/dxvgef/sessions")
	}
}

// DefaultCfg default config
func DefaultCfg() *Config {
	return &Config{CookieName: DefaultCookieName, Path: "/", MaxAge: DefaultMaxAge, HttpOnly: true, Secure: false}
}

// Context Session handle func
// return Session Item
func Context(w http.ResponseWriter, r *http.Request) (Session, error) {
	//防止处理时，进入另外的请求
	cookie, err := r.Cookie(_Cfg.CookieName)
	if err != nil || len(cookie.Value) <= 0 {
		item := recordCookie(w, cookie)
		return item, nil
	}
	// 防止浏览器关闭重新打开抛异常
	sid, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return nil, err
	}
	session := _Store.(*MemoryStorage).sessions[sid]
	if session == nil {
		item := recordCookie(w, cookie)
		return item, nil
	}
	return _Store.(*MemoryStorage).sessions[cookie.Value], nil
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
		MaxAge:   _Cfg.MaxAge,
		Expires:  time.Now().Add(time.Duration(_Cfg.MaxAge)),
	}
	http.SetCookie(w, cookie) //设置到响应中
	item := newSessionItem(sid, _Cfg.MaxAge)
	_Store.(*MemoryStorage).sessions[sid] = item
	return item
}
