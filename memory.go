// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 9:08 PM - UTC/GMT+08:00

package session

import (
	"strings"
	"sync"
	"time"
)

type MemoryStore struct {
	// lock When parallel, ensure data independence, consistency and safety
	mx sync.Mutex
	// sid:key:data
	values map[string]map[string][]byte
}

func newMemoryStore() *MemoryStore {
	//ms:= &MemoryStore{values: make(map[string]map[string]interface{}, MemoryMaxSize)}
	//ms.values[""] = make(map[string]interface{},maxSize)
	return &MemoryStore{values: make(map[string]map[string][]byte, MemoryMaxSize)}
}

func (m *MemoryStore) Writer(id, key string, data interface{}) error {
	if _Cfg._st == Memory {
		m.mx.Lock()
		defer m.mx.Unlock()
		// check map pointer is exist
		if m.values[id] == nil {
			m.values[id] = make(map[string][]byte, maxSize)
		}
		serialize, err := Serialize(data)
		if err != nil {
			return err
		}
		m.values[id][key] = serialize
		//log.Printf("%p",m.values[id])
		//log.Println(m.values[id][key])
		return nil
	}
	return nil
}
func (m *MemoryStore) Reader(id, key string) ([]byte, error) {
	if _Cfg._st == Memory {
		return m.values[id][key], nil
	}
	return nil, nil
}

func (m *MemoryStore) Remove(id, key string) {
	m.mx.Lock()
	defer m.mx.Unlock()
	delete(m.values[id], key)
}

func (m *MemoryStore) clean(id string) {
	m.values[id] = make(map[string][]byte, maxSize)
}

func (m *MemoryStore) gc() {
	for s, _ := range m.values {
		// 检测session是否过期 过期了就清理内存
		if time.Now().UnixNano() >= ParseInt64(strings.Split(s, ":")[1]) {
			delete(m.values, s)
		}
	}
}
