// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/24 - 3:58 下午 - UTC/GMT+08:00

package session

import "strconv"

// ParseString  int64  parse to string
func ParseString(nano int64) string {
	return strconv.FormatInt(nano, 10)
}

// ParseInt64 string parse to int64
// ret: -1 is  parse field
func ParseInt64(nano string) int64 {
	parseInt, err := strconv.ParseInt(nano, 10, 64)
	if err != nil {
		return -1
	}
	return parseInt
}
