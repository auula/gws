// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 9:36 下午 - UTC/GMT+08:00

package session

import (
	"log"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	_Cfg = &Config{MaxAge: ExpireTime}
	m.Run()
}

// go test -v -run=Test_newMemoryStore/test01
func Test_newMemoryStore(t *testing.T) {
	tests := []struct {
		name string
		want *MemoryStore
	}{
		{"test01", newMemoryStore()},
	}
	for _, tt := range tests {
		nano := time.Now().Add(time.Duration(_Cfg.MaxAge) * time.Second).UnixNano()
		formatInt := strconv.FormatInt(nano, 10)
		t.Log(formatInt)
		t.Run(tt.name, func(t *testing.T) {
			gm := make(map[string][]byte, 10)
			gm["t1"] = []byte("v1")
			wm := make(map[string][]byte, 5)
			wm["t2"] = []byte("v2")
			store := newMemoryStore()
			store.values["gotMap"] = gm
			tt.want.values["wantMap"] = wm
			t.Log(string(store.values["gotMap"]["t1"]) == string(tt.want.values["wantMap"]["t2"]))
			if got := store.values; !reflect.DeepEqual(got, tt.want.values) {
				t.Errorf("newMemoryStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStore_Writer(t *testing.T) {
	store := newMemoryStore()
	store.Writer("2006010213140506", "key1", "v1")
	t.Logf("%p", store.values["2006010213140506"])
	store.Writer("2006010213140506", "key2", "v2")
	t.Logf("%p", store.values["2006010213140506"])
	reader, _ := store.Reader("2006010213140506", "key1")
	log.Println(reader)
	log.Println(string(reader))
}
