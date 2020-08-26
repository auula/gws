// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/26 - 6:52 PM - UTC/GMT+08:00

package session

import "time"

var _GarbageList []*garbage

// 垃圾回收
type garbage struct {
	ID  string     // session ID
	Exp *time.Time // session expire
}

func remove(index int, gb []*garbage) []*garbage {
	gb = append(gb[:index], gb[index+1:]...)
	return gb
}

func containID(id string) bool {
	for _, g := range _GarbageList {
		return g.ID == id
	}
	return false
}

func gcPut(gb *garbage) {
	_GarbageList = append(_GarbageList, gb)
}
