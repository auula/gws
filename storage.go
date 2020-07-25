// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/7/25 - 12:35 PM

package session

// StoreType Custom Type
type StoreType uint8

const (
	MemoryType StoreType = iota
	RedisType
	DatabaseType
)

// session存储方式接口，可以存储在内存，数据库或者文件
// 分别实现该接口即可
// 如存入数据库的CRUD操作
type Storage interface {
	//初始化一个session，id根据需要生成后传入
	InitSession(sid string, maxAge int64) (Session, error)
	//根据sid，获得当前session
	SetSession(session Session) error
	//销毁session
	DestroySession(sid string) error
	//回收
	GCSession()
}
