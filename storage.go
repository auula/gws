// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 8:51 PM - UTC/GMT+08:00

package session

import "github.com/go-redis/redis"

// Storage is session store standard
type Storage interface {
	Writer(id, key string, data interface{}) error
	Reader(id, key string) ([]byte, error)
	Remove(id, key string)
	clean(id string)
}

// Value is unite session data value
type Value struct {
	Key   string
	Value string
	Error error
	*redis.StringCmd
}
