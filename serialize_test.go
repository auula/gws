// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/24 - 2:31 下午 - UTC/GMT+08:00

package session

import (
	"reflect"
	"testing"
)

func TestSerialize(t *testing.T) {

	type Users struct {
		Username string
		Password string
	}
	user := Users{
		"USER",
		"123456",
	}
	type args struct {
		obj interface{}
	}
	serialize, _ := Serialize(user)
	t.Log(string(serialize))
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"successful", args{
			obj: user,
		}, serialize},
		{"error", args{
			obj: user,
		}, []byte("111")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Serialize(tt.args.obj)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Serialize() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestUnSerialize(t *testing.T) {
	type User struct {
		Name string
		Age  int8
	}
	user := User{Name: "Ding", Age: 21}
	serialize, _ := Serialize(user)
	t.Log(serialize)
	var u User
	DeSerialize(serialize, &u)
	t.Log(u)
}

func Test_randomID(t *testing.T) {
	t.Log(string(randomID(16)))
	t.Log(string(Random(64, RuleKindAll)))
}
