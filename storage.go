// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/7/25 - 12:35 PM

package session

// Store  interface
type Store interface {
	// MemoryType
	// RedisType
	// DatabaseType
	Parse() (map[string]string, error)
}

// const (
// 	MemoryType StoreType = iota
// 	RedisType
// 	DatabaseType
// )
