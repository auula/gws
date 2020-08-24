// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 9:10 PM - UTC/GMT+08:00

package session

import (
	"net/http"
	"time"
)

// Session Unite struct
type Session struct {
	ID     string
	Cookie *http.Cookie
	Expire time.Time
}

// Capture return request session object
func Capture(writer http.ResponseWriter, request *http.Request) *Session {
	return nil
}

// Get get session data by key
func (s *Session) Get(key string) ([]byte, error) {
	if key == "" || len(key) <= 0 {
		return nil, ErrorKeyNotExist
	}
	//var result Value
	//result.Key = key
	b, err := _Store.Reader(s.parseID(), key)
	if err != nil {
		return nil, err
	}
	//result.Value = b
	return b, nil
}

// Set set session data by key
func (s *Session) Set(key string, data interface{}) error {
	if key == "" || len(key) <= 0 {
		return ErrorKeyFormat
	}
	return _Store.Writer(s.parseID(), key, data)
}

func (s *Session) parseID() (tmpId string) {
	switch _Cfg._st {
	case Memory:
		// 特殊格式sessionID 方便内存gc进行解析回收标识符
		tmpId = s.ID + ":" + ParseString(s.Expire.UnixNano())
	case Redis:
		tmpId = s.ID
	}
	return
}
