// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/7/24 - 07:47 PM

package session

import "sync"

// Version  for Session package
const Version = "0.0.1"

// Session Operation interface, Sesion operation of different storage methods is different,
//and the implementation is also different.
type Session interface {
    Set(key, value interface{})
    Get(key interface{}) interface{}
    Remove(key interface{}) error
    GetId() string
}

// SessionFromMemory implement
type SessionFromMemory struct {
    sid              string                      //唯一标示
    lock             sync.Mutex                  //一把互斥锁
    lastAccessedTime time.Time                   //最后访问时间
    maxAge           int64                       //超时时间
    data             map[interface{}]interface{} //主数据
}