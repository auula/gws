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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := rdb.Ping(ctx).Result()
	t.Log(result, err)

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
	rdb.Set(ctx, id, bytes, time.Second*10)

	s2 := new(Session)
	b, err := rdb.Get(ctx, id).Bytes()
	if err != nil {
		t.Log(err)
	}
	decoder(b, s2)
	t.Log(s2)
}
