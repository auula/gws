// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 8:51 PM - UTC/GMT+08:00

package session

import (
	"context"
	"time"
)

// Storage is session store standard
type Storage interface {
	Writer(ctx context.Context, key string, data interface{}) error
	Reader(id, key string) ([]byte, error)
	Remove(id, key string)
	Clean(id string)
}

// 后续版本会更新安全 1.解决了session会话超时时间伪造问题
// 这个是浏览器cookie同意保存的value值
// 把这个value加密然后响应给浏览器
// 服务器用秘钥解码就知道到期时间和数据了，这样客户端别人就伪造不了请求了
// Value is unite session data value
type Value struct {
	ID     string    `json:"id"`
	Expire time.Time `json:"expire"`
}
