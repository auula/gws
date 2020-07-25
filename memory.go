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

func newMemoryStore() *MemoryStore {
	return &MemoryStore{sessions: make(map[string]Session, 0)}
}

type Memory struct {
	SID              string                      // unique id
	Safe             sync.Mutex                  // mutex lock
	LastAccessedTime time.Time                   // last visit time
	MaxAge           int64                       // over time
	Data             map[interface{}]interface{} // save data
}

//实例化
func newMemory() *Memory {
	return &Memory{
		Data:   make(map[interface{}]interface{}),
		MaxAge: 60 * 30, //默认30分钟
	}
}

//同一个会话均可调用，进行设置，改操作必须拥有排斥锁
func (si *Memory) Set(key, value interface{}) {
	si.Safe.Lock()
	defer si.Safe.Unlock()
	si.Data[key] = value
}

func (si *Memory) Get(key interface{}) interface{} {
	if value := si.Data[key]; value != nil {
		return value
	}
	return nil
}
func (si *Memory) Remove(key interface{}) error {
	if value := si.Data[key]; value != nil {
		delete(si.Data, key)
	}
	return nil
}
func (si *Memory) GetId() string {
	return si.SID
}

//初始换会话session，这个结构体操作实现Session接口
func (fm *MemoryStore) InitSession(sid string, maxAge int64) (Session, error) {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	newSession := newMemory()
	newSession.SID = sid
	if maxAge != 0 {
		newSession.MaxAge = maxAge
	}
	newSession.LastAccessedTime = time.Now()

	fm.sessions[sid] = newSession //内存管理map
	return newSession, nil
}

//设置
func (fm *MemoryStore) SetSession(session Session) error {
	fm.sessions[session.GetId()] = session
	return nil
}

//销毁session
func (fm *MemoryStore) DestroySession(sid string) error {
	if _, ok := fm.sessions[sid]; ok {
		delete(fm.sessions, sid)
		return nil
	}
	return nil
}

//监判超时
func (fm *MemoryStore) GCSession() {

	sessions := fm.sessions

	//fmt.Println("gc session")

	if len(sessions) < 1 {
		return
	}

	//fmt.Println("current active sessions ", sessions)

	for k, v := range sessions {
		t := (v.(*Memory).LastAccessedTime.Unix()) + (v.(*Memory).MaxAge)

		if t < time.Now().Unix() { //超时了

			delete(fm.sessions, k)
		}
	}

}
