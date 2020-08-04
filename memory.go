// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/4 - 10:43 下午

package session

import "sync"

type MemoryStorage struct {
	//由于session包含所有的请求
	//并行时，保证数据独立、一致、安全
	lock     sync.Mutex //互斥锁
	sessions map[string]Session
	cfg      *Config
}

func NewMemoryStore() *MemoryStorage {
	return &MemoryStorage{sessions: make(map[string]Session, 128*maxSize)}
}
