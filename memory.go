// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 9:08 PM - UTC/GMT+08:00

package go_session

import "sync"

type MemoryStore struct {
	// lock When parallel, ensure data independence, consistency and safety
	mx sync.Mutex
	// sid:key:data
	values map[string]map[string]interface{}
}

func newMemoryStore() *MemoryStore {
	return &MemoryStore{values: make(map[string]map[string]interface{}, MemoryMaxSize)}
}
