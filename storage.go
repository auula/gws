// MIT License

// Copyright (c) 2022 Leon Ding

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

package gws

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// Storage global session data store interface.
// You can customize the storage medium by implementing this interface.
type Storage interface {
	Read(s *Session) (err error)
	Write(s *Session) (err error)
	Remove(s *Session) (err error)
}

// RamStore Local memory storage.
type RamStore struct {
	rw    *sync.RWMutex
	store map[string]*Session
}

// NewRAM return local memory storage.
func NewRAM() *RamStore {
	s := &RamStore{
		store: make(map[string]*Session),
		rw:    new(sync.RWMutex),
	}
	go s.gc()
	return s
}

func (ram *RamStore) Read(s *Session) (err error) {
	ram.rw.RLock()
	defer func() {
		if s.rw == nil {
			s.rw = new(sync.RWMutex)
		}
		ram.rw.RUnlock()
	}()
	if session, ok := ram.store[s.id]; ok {
		s.Values = session.Values
		s.CreateTime = session.CreateTime
		s.ExpireTime = session.ExpireTime
		return nil
	}
	debug.trace(s)
	return ErrSessionNoData
}

func (ram *RamStore) Write(s *Session) (err error) {
	ram.rw.Lock()
	defer ram.rw.Unlock()
	ram.store[s.id] = s
	debug.trace(s)
	return nil
}

func (ram *RamStore) Remove(s *Session) (err error) {
	ram.rw.Lock()
	defer ram.rw.Unlock()
	delete(ram.store, s.id)
	debug.trace(s)
	return nil
}

// gc is ram garbage collection.
func (ram *RamStore) gc() {
	for {
		// lifetime minute garbage collection.
		time.Sleep(lifeTime)
		debug.trace("gc running...")
		for _, session := range ram.store {
			if session.Expired() {
				ram.rw.Lock()
				delete(ram.store, session.id)
				ram.rw.Unlock()
			}
		}
	}
}

// RdsStore remote redis server storage.
type RdsStore struct {
	rw    *sync.RWMutex
	store *redis.Client
}

// NewRds return redis server storage.
func NewRds() *RdsStore {
	return &RdsStore{
		rw: new(sync.RWMutex),
		store: redis.NewClient(&redis.Options{
			Addr:     globalConfig.Address,
			Password: globalConfig.Password,
			DB:       int(globalConfig.Index),
			PoolSize: int(globalConfig.PoolSize),
		}),
	}
}

func (rds *RdsStore) Read(s *Session) (err error) {
	timeout, cancelFunc := timeoutCtx()
	rds.rw.RLock()
	defer func() {
		cancelFunc()
		// 防止读取的锁是空的，并且锁值不能共享
		// fix bug:
		// https://deepsource.io/gh/auula/gws/run/5b13c99b-9101-4e4f-8197-acfd730c28a0/go/VET-V0008
		if s.rw == nil {
			s.rw = new(sync.RWMutex)
		}
		rds.rw.RUnlock()
	}()
	var val []byte
	if val, err = rds.store.Get(timeout, formatPrefix(s.id)).Bytes(); err != nil {
		return err
	}
	debug.trace(val)
	return json.Unmarshal(val, s)
}

func (rds *RdsStore) Write(s *Session) (err error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}
	timeout, cancelFunc := timeoutCtx()
	rds.rw.Lock()
	defer func() {
		cancelFunc()
		rds.rw.Unlock()
	}()
	debug.trace(s)
	return rds.store.Set(timeout, formatPrefix(s.id), bytes, expire(s.ExpireTime)).Err()
}

func (rds *RdsStore) Remove(s *Session) (err error) {
	timeout, cancelFunc := timeoutCtx()
	rds.rw.Lock()
	defer func() {
		cancelFunc()
		rds.rw.Unlock()
	}()
	debug.trace(s)
	return rds.store.Del(timeout, formatPrefix(s.id)).Err()
}

func formatPrefix(sid string) string {
	return fmt.Sprintf("%s:%s", globalConfig.Prefix, sid)
}

func timeoutCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
}

func expire(t time.Time) time.Duration {
	return time.Until(t)
}
