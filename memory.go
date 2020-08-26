// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 9:08 PM - UTC/GMT+08:00

package session

import (
	"context"
	"errors"
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
	if ctx.Value(contextValue) == nil {
		return errors.New("context get value failed")
	}

	cv := ctx.Value(contextValue).(map[string]interface{})
	id := cv[contextValueID].(string)

	// 防止在写的过程中出现 gc回收了之前的id 程序找不到内存抛异常
	if m.values[id] == nil {
		m.values[id] = newMSessionItem(id, int(_Cfg.MaxAge))
	}

	m.values[id].Data[cv[contextValueKey].(string)] = cv[contextValueData]

	//log.Printf("%p",m.values[id])
	//log.Println(m.values[id][key])
	return nil
}

// Reader 读取数据 通过id和key
func (m *MemoryStore) Reader(ctx context.Context) ([]byte, error) {
	cv := ctx.Value(contextValue).(map[string]interface{})
	if len(cv[contextValueID].(string)) <= 0 {
		return nil, errors.New("context get uuid failed")
	}
	session := m.values[cv[contextValueID].(string)]
	if session == nil {
		return nil, errors.New("session id not  exist or not data")
	}
	return Serialize(session.Data[cv[contextValueKey].(string)])
}

// Remove 通过id和key移除数据
func (m *MemoryStore) Remove(ctx context.Context) {
	cv := ctx.Value(contextValue).(map[string]interface{})
	delete(m.values[cv[contextValueID].(string)].Data, cv[contextValueKey].(string))
}

// Remove 通过id和key移除数据
func (m *MemoryStore) Clean(ctx context.Context) {
	cv := ctx.Value(contextValue).(map[string]interface{})
	m.values[cv[contextValueID].(string)].Data = make(map[string]interface{}, maxSize)
}

// gc GarbageCollection
func (m *MemoryStore) gc() {
	// 10 分钟进行一次垃圾收集
	for {
		time.Sleep(10 * 60 * time.Second)
		sessions := m.values
		//if len(sessions) <= MemoryMaxSize/2 {
		//	continue
		//}
		for _, v := range sessions {
			if time.Now().UnixNano() >= v.Expires.UnixNano() { //超时了
				//fmt.Println("ID:", v.ID, "被GC清理了")
				delete(m.values, v.ID)
			}
		}
	}
}
