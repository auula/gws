// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 9:36 下午 - UTC/GMT+08:00

package go_session

import (
	"reflect"
	"testing"
)

// go test -v -run=Test_newMemoryStore/test01
func Test_newMemoryStore(t *testing.T) {
	tests := []struct {
		name string
		want *MemoryStore
	}{
		{"test01", newMemoryStore()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gm := make(map[string]interface{}, 10)
			gm["t1"] = "v1"
			wm := make(map[string]interface{}, 5)
			wm["t2"] = "v2"
			store := newMemoryStore()
			store.values["gotMap"] = gm
			tt.want.values["wantMap"] = wm
			t.Log(store.values["gotMap"]["t1"] == tt.want.values["wantMap"]["t2"])
			if got := store.values; !reflect.DeepEqual(got, tt.want.values) {
				t.Errorf("newMemoryStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
