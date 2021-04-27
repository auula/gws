// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2021/3/27 - 6:46 下午 - UTC/GMT+08:00

package sessionx

import (
	"sync"
	"testing"
)

func TestSessionLock(t *testing.T) {
	var count = 0
	var wait sync.WaitGroup
	wait.Add(3)
	// 并发写测试
	go lock["W"](func() {
		count = count + 1
		wait.Done()
	})
	// 并发读测试
	go lock["R"](func() {
		t.Log(count)
	})
	go lock["W"](func() {
		count = count + 1
		wait.Done()
	})
	go lock["W"](func() {
		count = count + 1
		wait.Done()
	})
	wait.Wait()
	t.Log(count)
}
