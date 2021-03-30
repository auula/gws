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
	"time"

	"github.com/go-redis/redis/v8"
)

// type task func() error

// redisStore redis storage implement
type redisStore struct {
	// channel  chan task
	sessions *redis.Client
}

// Read: read data
func (rs *redisStore) Read(s *Session) error {

	// set connection timeout
	timeout, cancelFunc := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancelFunc()

	sid := fmt.Sprintf("%s:%s", mgr.cfg.RedisKeyPrefix, s.ID)
	bytes, err := rs.sessions.Get(timeout, sid).Bytes()
	if err != nil {
		return err
	}
	// Expire timeout context
	timeout, cancelFunc = context.WithTimeout(context.Background(), 8*time.Second)
	defer cancelFunc()

	if err := rs.sessions.Expire(timeout, sid, mgr.cfg.TimeOut).Err(); err != nil {
		return err
	}
	if err := decoder(bytes, s); err != nil {
		return err
	}
	// log.Println("redis read:", s)
	return nil
}

// Create: create session data
func (rs *redisStore) Create(s *Session) error {
	return rs.setValue(s)
}

// Update: updated session data
func (rs *redisStore) Update(s *Session) error {
	return rs.setValue(s)
}

// Remove: remove session data
func (rs *redisStore) Remove(s *Session) error {
	// set connection timeout
	timeout, cancelFunc := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancelFunc()
	return rs.sessions.Del(timeout, fmt.Sprintf("%s:%s", mgr.cfg.RedisKeyPrefix, s.ID)).Err()
}

// setValue: set session data value
func (rs *redisStore) setValue(s *Session) error {
	bytes, err := encoder(s)
	if err != nil {
		return err
	}
	// set connection timeout
	timeout, cancelFunc := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancelFunc()

	return rs.sessions.Set(timeout, fmt.Sprintf("%s:%s", mgr.cfg.RedisKeyPrefix, s.ID), bytes, mgr.cfg.TimeOut).Err()
}

//func (rs *redisStore) Do()  {
//
//}
