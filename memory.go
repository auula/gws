// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 9:08 PM - UTC/GMT+08:00

package session

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MemoryStore 内存存储实现
type MemoryStore struct {
	// lock When parallel, ensure data independence, consistency and safety
	mx sync.Mutex
	// sid:key:data save serialize data
	values map[string]*MemorySession
}

// newMemoryStore 创建一个内存存储 开辟内存
func newMemoryStore() *MemoryStore {
	ms := &MemoryStore{values: make(map[string]*MemorySession, MemoryMaxSize)}
	//ms.values[""] = make(map[string]interface{},maxSize)
	_GarbageList = make([]*garbage, 0, MemoryMaxSize)
	go ms.gc()
	return ms
}

// Writer 写入数据方法
func (m *MemoryStore) Writer(ctx context.Context) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	// check map pointer is exist
	cv := ctx.Value(contextValue).(map[string]interface{})
	id := cv[contextValueID].(string)
	if m.values[id] == nil {
		// 方便后面进行gc()
		gcPut(&garbage{ID: id, Exp: cv[contextValueExpire].(*time.Time)})
	}
	serialize, err := Serialize(cv[contextValueData].(*))
	if err != nil {
		return err
	}
	m.values[id][key] = serialize
	//log.Printf("%p",m.values[id])
	//log.Println(m.values[id][key])
	return nil
}

// Reader 读取数据 通过id和key
func (m *MemoryStore) Reader(id, key string) ([]byte, error) {
	return m.values[id][key], nil
}

// Remove 通过id和key移除数据
func (m *MemoryStore) Remove(id, key string) {
	delete(m.values[id], key)

}

// Clean 通过id清空data
func (m *MemoryStore) Clean(id string) {
	m.values[id] = make(map[string][]byte, maxSize)
}

// gc GarbageCollection
func (m *MemoryStore) gc() {
	var index int
	for {
		time.Sleep(10 * time.Second)
		for i, g := range _GarbageList {
			index = i
			fmt.Println(g.ID, g.Exp.UnixNano())
			if time.Now().UnixNano() >= g.Exp.UnixNano() {
				delete(m.values, g.ID)
			}
		}
		if len(_GarbageList) > 0 {
			// 移除垃圾堆里面的session
			_GarbageList = remove(index, _GarbageList)
		}
	}
}
