// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 9:29 PM - UTC/GMT+08:00

package session

import (
	"errors"
)

type StoreType int8

const (
	Memory StoreType = iota
	Redis
	maxSize           = 16
	DefaultCookieName = "go_session_key"
	DefaultMaxAge     = 30 * 60 // 30 min
	ExpireTime        = 1800
	RedisMaxSize      = 16
	MemoryMaxSize     = 16 * 1024
	RedisPrefix       = "go:session:"

	// RuleKindNumber is Number
	RuleKindNumber = iota
	// RuleKindLower is letter lower
	RuleKindLower
	// RuleKindUpper is letter upper
	RuleKindUpper
	// RuleKindAll is all rule
	RuleKindAll
)

var (
	_Cfg             *Config
	_Store           Storage
	ErrorKeyFormat   = errors.New("set session data failed,key format error")
	ErrorKeyNotExist = errors.New("get session data failed,key not exist")
	ErrorSetValue    = errors.New("set session data failed,redis save failed")
)

// Config param
type Config struct {
	// cookie参数
	CookieName     string // sessionID的cookie键名
	Domain         string // sessionID的cookie作用域名
	Path           string // sessionID的cookie作用路径
	MaxAge         int64  // 最大生命周期（秒）
	HttpOnly       bool   // 仅用于http（无法被js读取）
	Secure         bool   // 启用https
	EncryptedKey   string // sessionID值加密的密钥
	RedisAddr      string // redis地址
	RedisPassword  string // redis密码
	RedisKeyPrefix string // redis键名前缀
	//IdleTime                  time.Duration // 空闲生命周期
	RedisDB int // redis数据库
	//DisableAutoUpdateIdleTime bool          // 禁止自动更新空闲时间
	_st StoreType
}

// DefaultCfg default config
func DefaultCfg() *Config {
	return &Config{CookieName: DefaultCookieName, Path: "/", MaxAge: 60, HttpOnly: true, Secure: false}
}

// ReloadCfg reload config
func ReloadCfg(config *Config) {
	_Cfg = config
}
