// MIT License

// Copyright (c) 2022 Leon Ding

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package gws

import (
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewSession(t *testing.T) {
	nowTime := time.Now()
	uuid := uuid73()

	session := &Session{
		session{
			id:         uuid,
			rw:         sync.RWMutex{},
			Values:     make(Values),
			CreateTime: nowTime,
			ExpireTime: nowTime.Add(lifeTime),
		},
	}

	tests := []struct {
		name string
		want *Session
	}{
		{"successful", session},
		//{"fail",NewSession()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := session; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_Expired(t *testing.T) {
	nowTime := time.Now()
	uuid := uuid73()

	type fields struct {
		session session
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"successful", fields{
			session: session{
				id:         uuid,
				rw:         sync.RWMutex{},
				Values:     make(Values),
				CreateTime: nowTime,
				ExpireTime: nowTime.Add(lifeTime),
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				session: tt.fields.session,
			}
			if got := s.Expired(); got != tt.want {
				t.Errorf("Expired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCookie(t *testing.T) {

	Open(DefaultRAMOptions)

	tests := []struct {
		name string
		want *http.Cookie
	}{
		{"successful", NewCookie()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCookie(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCookie() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreFactory(t *testing.T) {

	store := NewRAM()
	StoreFactory(NewOptions(), store)

	type args struct {
		opt   Options
		store Storage
	}
	tests := []struct {
		name string
		args args
	}{
		{"successful", args{
			NewOptions(),
			store,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if &tt.args.store == &globalStore {
				t.Error("store init fail")
			}
		})
	}
}
