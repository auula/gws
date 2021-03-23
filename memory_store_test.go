package sessionx

import (
	"github.com/google/uuid"
	"net/http"
	"testing"
	"time"
)

var (
	_testCfg = &Configs{
		EncryptedKey:   "0123456789012345",
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
)

func TestCreate(t *testing.T) {
	m := new(memoryStore)
	s := new(Session)
	s.ID = "20210320"
	m.Create(s)
	t.Log("session = ", s)
}

func TestReader(t *testing.T) {
	m := new(memoryStore)
	s := new(Session)
	s.ID = "20210320"
	m.Create(s)
	err := m.Reader(s)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("session = ", s)
}

func TestDelete(t *testing.T) {
	m := new(memoryStore)
	s := new(Session)
	s.ID = "20210320"
	m.Create(s)
	t.Log("session = ", s)
	m.Delete(s)
	err := m.Reader(s)
	if err != nil {
		t.Error(err.Error())
	}

}

func TestUpdated(t *testing.T) {
	m := new(memoryStore)
	s := new(Session)
	s.ID = "20210320"
	m.Create(s)
	t.Log("session = ", s)
	v := make(map[string]interface{})
	v["v"] = "test"
	m.Update(&Session{ID: "20210320", Data: v, Expires: time.Now()})
	err := m.Reader(s)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("session = ", s)

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

	New(M, _testCfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := new(Session)
		s.ID = uuid.New().String()
		s.Data = make(map[string]interface{}, 8)
		s.Expires = time.Now().Add(mgr.cfg.TimeOut)
		_ = mgr.store.Update(s)
	}
}
