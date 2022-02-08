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

//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"net/http"

	"github.com/auula/gws"
)

func init() {
	gws.Debug(false)
	gws.StoreFactory(gws.NewOptions(), &FileStore{})
}

type FileStore struct{}

func (fs FileStore) Read(s *gws.Session) (err error) {
	panic("implement me")
}

func (fs FileStore) Write(s *gws.Session) (err error) {
	panic("implement me")
}

func (fs FileStore) Remove(s *gws.Session) (err error) {
	panic("implement me")
}

func main() {

	http.HandleFunc("/panic", func(writer http.ResponseWriter, request *http.Request) {

		session, _ := gws.GetSession(writer, request)
		session.Values["foo"] = "bar"
		session.Sync()

		fmt.Fprintln(writer, "set value successful.")
	})

	_ = http.ListenAndServe(":8080", nil)

}
