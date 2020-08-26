// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 9:08 PM - UTC/GMT+08:00

package session

import (
	"context"
	"errors"
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
	if ctx.Value(contextValue) == nil {
		return errors.New("context get value failed")
	}

	cv := ctx.Value(contextValue).(map[string]interface{})
	id := cv[contextValueID].(string)

	if m.values[id] == nil {
		// 往gc切片里面添加过期项 方便后面进行gc()
		gcPut(&garbage{
			ID:  id,
			Exp: cv[contextValueExpire].(*time.Time)})
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
	fmt.Println("GC START")
	var index int
	for {
		time.Sleep(10 * time.Second)
		for i, g := range _GarbageList {
			index = i
			fmt.Println("GC START 1")
			fmt.Println(g.ID, g.Exp.UnixNano())
			if time.Now().UnixNano() >= g.Exp.UnixNano() {
				fmt.Println("GC START 2")
				fmt.Println("销毁:", g.ID)
				delete(m.values, g.ID)
			}
		}
		if len(_GarbageList) > 0 {
			// 移除垃圾堆里面的session
			_GarbageList = remove(index, _GarbageList)
		}
	}
}
