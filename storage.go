// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/4 - 10:44 PM

package session

// Store standard
type Store interface {
	// Garbage Collection
	GC()
}
