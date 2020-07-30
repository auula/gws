// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/7/25 - 12:36 PM

package session

import (
	"sync"
	"time"
)

type MemoryStore struct {
	//由于session包含所有的请求
	//并行时，保证数据独立、一致、安全
	lock     sync.Mutex //互斥锁
	sessions map[string]Session
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{sessions: make(map[string]Session, 0)}
}

type MemoryItem struct {
	SID              string                      // unique id
	Safe             sync.Mutex                  // mutex lock
	LastAccessedTime time.Time                   // last visit time
	MaxAge           int64                       // over time
	Data             map[interface{}]interface{} // save data
}

//实例化
func newMemoryItem(id string) *MemoryItem {
	return &MemoryItem{
		Data:   make(map[interface{}]interface{}),
		MaxAge: 60 * 30, //默认30分钟
		SID:    id,
	}
}

//同一个会话均可调用，进行设置，改操作必须拥有排斥锁
func (si *MemoryItem) Set(key, value interface{}) {
	si.Safe.Lock()
	defer si.Safe.Unlock()
	si.Data[key] = value
}

func (si *MemoryItem) Get(key interface{}) interface{} {
	if value := si.Data[key]; value != nil {
		return value
	}
	return nil
}

func (si *MemoryItem) Remove(key interface{}) error {
	if value := si.Data[key]; value != nil {
		delete(si.Data, key)
	}
	return nil
}

func (si *MemoryItem) GetID() string {
	return si.SID
}

// Parse parameter
func (ms *MemoryStore) Parse() (map[string]string, error) {
	// Not to do
	return nil, nil
}

// GCSession 监判超时
func (ms *MemoryStore) GC() {
	sessions := ms.sessions
	//fmt.Println("gc session")
	if len(sessions) < 1 {
		return
	}
	//fmt.Println("current active sessions ", sessions)
	for k, v := range sessions {
		t := (v.(*MemoryItem).LastAccessedTime.Unix()) + (v.(*MemoryItem).MaxAge)
		if t < time.Now().Unix() { //超时了
			delete(ms.sessions, k)
		}
	}
}
