// MIT License

// Copyright (c) 2021 Jarvib Ding

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package sessionx

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
)

type storeType uint8

const (
	// memoryStore store type
	M storeType = iota
	// redis store type
	R
	SessionKey = "session-id"
)

// manager for session manager
type manager struct {
	cfg   *Configs
	store storage
}

func New(t storeType, cfg *Configs) {

	switch t {
	case M:
		// init memory storage
		m := new(memoryStore)
		go m.gc()
		mgr = &manager{cfg: cfg, store: m}

	case R:

		// parameter verify
		validate := validator.New()
		if err := validate.Struct(cfg); err != nil {
			panic(err.Error())
		}

		// init redis storage
		r := new(redisStore)
		r.sessions = redis.NewClient(&redis.Options{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword, // no password set
			DB:       cfg.RedisDB,       // use default DB
			PoolSize: int(cfg.PoolSize), // connection pool size
		})

		// test connection
		timeout, cancelFunc := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancelFunc()
		if err := r.sessions.Ping(timeout).Err(); err != nil {
			panic(err.Error())
		}
		mgr = &manager{cfg: cfg, store: r}

	default:
		panic("not implement store type")
	}
}
