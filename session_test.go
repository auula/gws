// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2021/3/27 - 6:46 下午 - UTC/GMT+08:00

package sessionx

import (
	"github.com/google/uuid"
	"net/http"
	"testing"
	"time"
)

var (
	// 测试配置
	tc = &Configs{
		TimeOut:        time.Minute * 30,
		RedisAddr:      "127.0.0.1:6379",
		RedisDB:        0,
		RedisPassword:  "redis.nosql",
		RedisKeyPrefix: SessionKey,

		Cookie: &http.Cookie{
			Name:     SessionKey,
			Path:     "/",
			Expires:  time.Now().Add(time.Minute * 30),
			Secure:   false,
			HttpOnly: true,
			MaxAge:   60 * 30,
		},
	}
	ms *memoryStore
)

func TestSession(t *testing.T) {
	ms = new(memoryStore)
	go ms.gc()
	mgr = &manager{cfg: _testCfg, store: ms}
	ss := new(Session)
	ss.ID = uuid.New().String()
	ss.Data = make(map[interface{}]interface{}, 8)
	ss.Cookie = tc.Cookie
	ss.Expires = time.Now().Add(_testCfg.TimeOut)
	ss._w = nil

	ss.Set("k", "v")
	t.Log(ss.Get("k"))
	ss.Remove("k")
	t.Log(ss.Get("k"))
}
