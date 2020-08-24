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
	// sid:key:data save serialize data
	values map[string]map[string][]byte
}

func newMemoryStore() *MemoryStore {
	ms := &MemoryStore{values: make(map[string]map[string][]byte, MemoryMaxSize)}
	//ms.values[""] = make(map[string]interface{},maxSize)
	go ms.gc()
	return ms
}

func (m *MemoryStore) Writer(id, key string, data interface{}) error {
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
func (m *MemoryStore) Reader(id, key string) ([]byte, error) {
	return m.values[id][key], nil
}

func (m *MemoryStore) Remove(id, key string) {
	delete(m.values[id], key)
}

func (m *MemoryStore) Clean(id string) {
	m.values[id] = make(map[string][]byte, maxSize)
}

func (m *MemoryStore) gc() {
	for {
		// 每30分钟进行一次垃圾清理  session过期的全部清理掉
		time.Sleep(30 * 60 * time.Second)
		if len(m.values) < 1 {
			continue
		}
		for s, _ := range m.values {
			if time.Now().UnixNano() >= ParseInt64(strings.Split(s, ":")[1]) {
				delete(m.values, s)
			}
		}
	}
}
