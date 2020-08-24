// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/24 - 2:15 下午 - UTC/GMT+08:00

package session

import (
	"bytes"
	"encoding/gob"
	"fmt"
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

// DeSerialize byte to object
func DeSerialize(byte []byte, obj interface{}) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(byte)
	if err != nil {
		return err
	}
	decoder := gob.NewDecoder(bytes.NewReader(buf.Bytes()))
	fmt.Println(decoder.Decode(&obj))
	return decoder.Decode(&obj)
}
