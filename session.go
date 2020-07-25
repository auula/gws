// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/7/24 - 07:47 PM

package session

// Version  for Session package
const Version = "0.0.1"

// Session Operation interface, Session operation of different storage methods is different,
// and the implementation is also different.
type Session interface {
	Set(key, value interface{})
	Get(key interface{}) interface{}
	Remove(key interface{}) error
	GetId() string
}

func Builder(storeType StoreType, sid string, maxAge int64) (Session, error) {
	switch storeType {
	case MemoryType:
		return newMemoryStore().InitSession(sid, maxAge)
	case RedisType:
		panic("not implement store type!")
	case DatabaseType:
		panic("not implement store type!")
	default:
		panic("not implement store type!")
	}
	//return nil,errors.New("build store type failed")
}
