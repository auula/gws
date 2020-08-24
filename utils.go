// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/24 - 3:58 下午 - UTC/GMT+08:00

package session

import (
	"math/rand"
	"strconv"
	"time"
)

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

func randomID(len int8) []byte {
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, 0, len)
	for i := 0; i < cap(result); i++ {
		result = append(result, byte(rand.Intn(26)+65))
	}
	return result
}

// Random is random string func
// ASCII CODE
func Random(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	// don't use underscores in Go names; var is_all should be isAll ！！！
	// 我要记住这些规范！！！
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			ikind = rand.Intn(3)
		}
		// https://iknow-pic.cdn.bcebos.com/aa18972bd40735fa8219094e92510fb30e24085f?x-bce-process=image/resize,m_lfit,w_600,h_800,limit_1
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
		// num:   0872394784259522
		// lower: impedfcjyeipdtuo
		// upper: IMPEDFCJYEIPDTUO
		// all:   Me9jEpT2WYMayeNZ
	}
	return result
}
