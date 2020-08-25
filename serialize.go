// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/24 - 2:15 下午 - UTC/GMT+08:00

package session

import (
	"encoding/json"
)

// Serialize object to serialize byte
// Serialize 只限于 interface{} 类型 转换 []byte使用
// 严禁把 string 类型 传入 obj 当做转换对象使用
// 大坑 大坑 大坑！！！！debug你的debug不出来的！！！！！
// 笔者友情提示！！！！！
// 2020-08-25 12:42:57 By: SDing
// 问题解析 : https://www.jianshu.com/p/f778206ac54c
func Serialize(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

// DeSerialize byte to object
func DeSerialize(byte []byte, obj interface{}) error {
	return json.Unmarshal(byte, obj)
}
