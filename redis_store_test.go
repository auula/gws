// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2021/3/29 - 6:29 下午 - UTC/GMT+08:00

package sessionx

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	rdb *redis.Client
	w   sync.WaitGroup
	mux sync.Mutex
	tx  = make(chan struct{}, 1)
)

func TestCoroutine(t *testing.T) {
	ctx := context.Background()
	rdb = redis.NewClient(
		&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "redis.nosql",
			DB:       0,
		},
	)
	w.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer w.Done()
			for {
				select {
				case <-tx:
					s := rdb.Get(ctx, "i").String()
					num, _ := strconv.Atoi(s)
					num = num + 1
					rdb.Set(ctx, "i", num, 10*time.Minute)
					releaseLock()
					break
				default:
					if upLock() {
						tx <- struct{}{}
					}
				}
			}

		}()
	}
	w.Wait()
	t.Log("i = ", rdb.Get(ctx, "i"))
}

// 如果lockChan中为空则阻塞
func upLock() bool {
	// 设置适当的过期时间，若出现断连情况，redis会将超时的锁删除，防止死锁
	result, err := rdb.SetNX(ctx, "lock", 1, time.Millisecond*1000*20).Result()
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("result",result)
	//fmt.Println(result)
	return result
}

// 重新填充lockChan
func releaseLock() {
	rdb.Del(ctx, "lock")
}

func Affairs(fn func()) {
	if upLock() {
		defer releaseLock()
		fn()
	}
	Affairs(fn)
}
