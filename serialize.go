// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/24 - 2:15 下午 - UTC/GMT+08:00

package session

import (
	"encoding/json"
)

// Serialize object to serialize byte
func Serialize(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

// DeSerialize byte to object
func DeSerialize(byte []byte, obj interface{}) error {
	return json.Unmarshal(byte, obj)
}
