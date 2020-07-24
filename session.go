// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/7/24 - 07:47 PM

package session

import (
	"sync"
	"time"
)

// Version  for Session package
const Version = "0.0.1"

// Session Operation interface, Sesion operation of different storage methods is different,
// and the implementation is also different.
type Session interface {
	Set(key, value interface{})
	Get(key interface{}) interface{}
	Remove(key interface{}) error
	GetId() string
}

// FromMemory struct data store in Redis
type FromMemory struct {
	sid              string                      // unique id
	lock             sync.Mutex                  // mutex lock
	lastAccessedTime time.Time                   // last visit time
	maxAge           int64                       // over time
	data             map[interface{}]interface{} // save data
}

// FromRedis struct  data store in Redis
type FromRedis struct {
	sid              string                      // unique id
	lock             sync.Mutex                  // mutex lock
	lastAccessedTime time.Time                   // last visit time
	maxAge           int64                       // over time
	data             map[interface{}]interface{} // save data
}
