// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/7/25 - 6:42 PM

package session

import (
	"reflect"
	"testing"
)

// go test -v -run=TestNewMgr/successful || error
func TestNewMgr(t *testing.T) {
	type args struct {
		storeType  StoreType
		cookieName string
		maxAge     int64
	}
	tests := []struct {
		name string
		args args
		want *Manager
	}{
		// TODO: Add test cases.
		{"successful", args{maxAge: 30 * 60, cookieName: "TestName", storeType: MemoryType}, New(MemoryType, "TestName", 30*60)},
		{"error", args{maxAge: 30 * 60, cookieName: "TestNamexxx", storeType: MemoryType}, New(MemoryType, "TestName", 30*60)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.storeType, tt.args.cookieName, tt.args.maxAge); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
