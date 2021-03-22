package sessionx

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"testing"
	"time"
)

func Test_redisStore(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redis.nosql", // no password set
		DB:       0,             // use default DB
		PoolSize: 100,           // 连接池大小
	})

	s := new(Session)
	id := uuid.New().String()
	s.ID = id
	s.Expires = time.Now()
	s.Data = make(map[string]interface{})
	s.Data["k"] = "v"
	bytes, err := encoder(s)
	if err != nil {
		t.Log(err)
	}
	rdb.Set(context.Background(), id, bytes, time.Second*10)

	s2 := new(Session)
	b, err := rdb.Get(context.Background(), id).Bytes()
	if err != nil {
		t.Log(err)
	}
	_ = decoder(b, s2)
	t.Log(s2)
}
