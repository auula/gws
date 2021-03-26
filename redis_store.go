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
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
)

var (
	ctx = context.Background()
)

// 关闭之后 redis有数据  但是获取2次就没有数据了
type redisStore struct {
	sync.Mutex
	sessions *redis.Client
}

func (rs *redisStore) Reader(s *Session) error {
	sid := fmt.Sprintf("%s:%s", mgr.cfg.RedisKeyPrefix, s.ID)
	rs.Lock()
	defer rs.Unlock()
	bytes, err := rs.sessions.Get(ctx, sid).Bytes()
	if err != nil {
		return err
	}
	if err := rs.sessions.Expire(ctx, sid, mgr.cfg.TimeOut).Err(); err != nil {
		return err
	}
	if err := decoder(bytes, s); err != nil {
		return err
	}
	log.Println("redis read:", s)
	return nil
}

func (rs *redisStore) Create(s *Session) error {
	rs.Lock()
	defer rs.Unlock()
	if s.Data == nil {
		s.Data = make(map[string]interface{}, 8)
	}
	return rs.setValue(s)
}

func (rs *redisStore) Update(s *Session) error {

	rs.Lock()
	defer rs.Unlock()

	if s.Data == nil {
		s.Data = make(map[string]interface{}, 8)
	}
	return rs.setValue(s)
}

func (rs *redisStore) Remove(s *Session, key string) error {
	rs.Lock()
	defer rs.Unlock()
	// delete it form memory
	if _, ok := s.Data[key]; ok {
		delete(s.Data, key)
	}
	return rs.setValue(s)
}

func (rs *redisStore) Delete(s *Session) error {
	rs.Lock()
	defer rs.Unlock()
	return rs.sessions.Del(ctx, fmt.Sprintf("%s:%s", mgr.cfg.RedisKeyPrefix, s.ID)).Err()
}

func (rs *redisStore) setValue(s *Session) error {
	bytes, err := encoder(s)
	if err != nil {
		return err
	}
	err = rs.sessions.Set(ctx, fmt.Sprintf("%s:%s", mgr.cfg.RedisKeyPrefix, s.ID), bytes, mgr.cfg.TimeOut).Err()
	return err
}
