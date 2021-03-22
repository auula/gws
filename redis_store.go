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
	"sync"
	"time"
)

var ()

type redisStore struct {
	ctx context.Context
	sync.Mutex
	sessions *redis.Client
	cancel   context.CancelFunc
}

func (rs *redisStore) Reader(s *Session) ([]byte, error) {
	rs.Lock()
	defer rs.Unlock()
	bytes, err := rs.sessions.Get(rs.ctx, fmt.Sprintf("%s:%s", mgr.cfg.RedisKeyPrefix, s.ID)).Bytes()
	if err != nil {
		return nil, err
	}
	err = decoder(bytes, s)
	if err != nil {
		return nil, err
	}
	return encoder(s)
}

func (rs *redisStore) Create(s *Session) ([]byte, error) {
	rs.Lock()
	defer rs.Unlock()
	if s.Data == nil {
		s.Data = make(map[string]interface{}, 8)
	}
	return rs.setValue(s)
}

func (rs *redisStore) Update(s *Session) ([]byte, error) {
	rs.Lock()
	defer rs.Unlock()
	s.Expires = time.Now().Add(mgr.cfg.SessionLifeTime)
	return rs.setValue(s)
}

func (rs *redisStore) Remove(s *Session, key string) error {
	rs.Lock()
	defer rs.Unlock()
	s.Expires = time.Now().Add(mgr.cfg.SessionLifeTime)
	// delete it form memory
	if _, ok := s.Data[key]; ok {
		delete(s.Data, key)
	}
	_, err := rs.setValue(s)
	return err
}

func (rs *redisStore) Delete(s *Session) error {
	rs.Lock()
	defer rs.Unlock()
	return rs.sessions.Del(rs.ctx, fmt.Sprintf("%s:%s", mgr.cfg.RedisKeyPrefix, s.ID)).Err()
}

func (rs *redisStore) setValue(s *Session) ([]byte, error) {
	bytes, err := encoder(s)
	if err != nil {
		return nil, err
	}
	err = rs.sessions.Set(rs.ctx, fmt.Sprintf("%s:%s", mgr.cfg.RedisKeyPrefix, s.ID), bytes, -1).Err()
	return nil, err
}
