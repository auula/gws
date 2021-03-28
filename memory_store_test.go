// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2021/3/26 - 11:48 下午 - UTC/GMT+08:00

package sessionx

import (
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	// 测试配置
	_testCfg = &Configs{
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
	m *memoryStore
	s *Session
)

func TestMain(t *testing.M) {
	m = new(memoryStore)
	go m.gc()
	mgr = &manager{cfg: _testCfg, store: m}

	s = new(Session)
	s.ID = uuid.New().String()
	s.Data = make(map[string]interface{}, 8)
	s.Cookie = _testCfg.Cookie
	s.Expires = time.Now().Add(_testCfg.TimeOut)
	t.Run()
}

func TestALL(t *testing.T) {

	m.Create(s)
	t.Log("Create session = ", s)

	v := make(map[string]interface{})
	v["v"] = "test"
	s.Data = v
	m.Update(s)
	err := m.Read(s)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Read session = ", s)

	m.Delete(s.ID)
	err = m.Read(s)
	if err != nil {
		t.Log("Delete session successful ")
	}

}

// https://my.oschina.net/solate/blog/3034188
func BenchmarkWrite(b *testing.B) {
	// 60s test
	//go test -bench=. -benchtime=60s -run=none
	//goos: windows
	//goarch: amd64
	//pkg: github.com/higker/sesssionx
	//cpu: Intel(R) Core(TM) i5-7200U CPU @ 2.50GHz
	//	BenchmarkWrite-4        95464014               810.1 ns/op
	//	PASS
	//	ok      github.com/higker/sesssionx     80.857s

	//go test -bench=. -run=none
	//goos: windows
	//goarch: amd64
	//pkg: github.com/higker/sesssionx
	//cpu: Intel(R) Core(TM) i5-7200U CPU @ 2.50GHz
	//	BenchmarkWrite-4         1569414               758.8 ns/op
	//	PASS
	//	ok      github.com/higker/sesssionx     3.664s

	b.Logf("系统:%s CPU核数:%d ", runtime.GOOS, runtime.NumCPU())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.ID = uuid.New().String()
		_ = m.Update(s)
	}
}
