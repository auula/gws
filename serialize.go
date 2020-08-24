// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/24 - 2:15 下午 - UTC/GMT+08:00

package session

import (
	"bytes"
	"encoding/gob"
)

// Serialize object to serialize byte
func Serialize(obj interface{}) ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(obj)
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

// UnSerialize byte to object
func UnSerialize(byte []byte, obj interface{}) {
	decoder := gob.NewDecoder(bytes.NewReader(byte))
	decoder.Decode(&obj)
}
