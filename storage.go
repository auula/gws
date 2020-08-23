// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 8:51 PM - UTC/GMT+08:00

package go_session

type Storage interface {
	Writer(key string, data interface{})
	Reader(key string) (interface{}, error)
	Remove(key string)
	Clean(key string)
}
