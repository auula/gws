// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/4 - 10:43 PM

package session

import (
	"fmt"
	"sync"
	"time"
)

type MemoryStorage struct {
	//由于session包含所有的请求
	//并行时，保证数据独立、一致、安全
	lock     sync.Mutex //互斥锁
	sessions map[string]Session
}

func newMemoryStore() *MemoryStorage {
	return &MemoryStorage{sessions: make(map[string]Session, 128*maxSize), cfg: config}
}

func (ms *MemoryStorage) GC() {
	// 10 分钟进行一次垃圾收集
	for {
		time.Sleep(10 * time.Second)
		sessions := ms.sessions
		if len(sessions) < 1 {
			continue
		}
		for k, v := range sessions {
			fmt.Println(time.Now())
			fmt.Println(v.(*Item).expires)
			if time.Now().Unix() >= v.(*Item).expires.Unix() { //超时了
				fmt.Println("ID:", k, "被GC清理了")
				delete(ms.sessions, k)
			}
		}
	}
}

// Item a session  item
type Item struct {
	sID     string                      // unique id
	safe    sync.Mutex                  // mutex lock
	expires time.Time                   // Expires time
	data    map[interface{}]interface{} // save data
}

//实例化
func newSessionItem(id string, maxAge int) *Item {
	return &Item{
		data:    make(map[interface{}]interface{}, maxSize),
		sID:     id,
		expires: time.Now().Add(time.Duration(maxAge) * time.Second),
	}
}

//同一个会话均可调用，进行设置，改操作必须拥有排斥锁
func (si *Item) Set(key, value interface{}) {
	si.safe.Lock()
	defer si.safe.Unlock()
	si.data[key] = value
}

func (si *Item) Get(key interface{}) interface{} {
	if value := si.data[key]; value != nil {
		return value
	}
	return nil
}

func (si *Item) Remove(key interface{}) {
	if value := si.data[key]; value != nil {
		delete(si.data, key)
	}
}

func (si *Item) ID() string {
	return si.sID
}
func (si *Item) Clear() {
	si.data = make(map[interface{}]interface{}, maxSize)
}
