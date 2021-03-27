// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2021/3/27 - 6:46 下午 - UTC/GMT+08:00

package sessionx

import (
	"sync"
	"testing"
)

var (
	rwm  sync.RWMutex
	lock = map[string]method{
		"W": func(f func()) {
			rwm.Lock()
			defer rwm.Unlock()
			f()
		},
		"R": func(f func()) {
			rwm.RLock()
			defer rwm.RUnlock()
			f()
		},
	}
)

func TestSessionLock(t *testing.T) {
	var count = 0
	var wait sync.WaitGroup
	wait.Add(3)
	go lock["W"](func() {
		count += 1
		t.Log(count)
	})
	go lock["W"](func() {
		count += 1
		t.Log(count)
	})
	go lock["W"](func() {
		count += 1
		t.Log(count)
	})
	wait.Done()
	t.Log(count)
}
