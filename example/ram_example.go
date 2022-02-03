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

//+build ignore

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/auula/gws"
)

func init() {
	gws.Debug(true)

	// var opt gws.RAMOption
	// opt.Domain = "www.ibyte.me"
	// gws.Open(opt)

	// gws.Open(gws.NewRDSOptions("", 1999, ""))
	// gws.Open(gws.NewOptions())
	// gws.Open(gws.NewOptions(gws.Domain(""), gws.CookieName("")))

	gws.Open(gws.DefaultRAMOptions)

}

type UserInfo struct {
	UserName string `json:"user_name,omitempty"`
	Email    string `json:"email,omitempty"`
	Age      uint8  `json:"age,omitempty"`
}

func main() {
	http.HandleFunc("/set", func(writer http.ResponseWriter, request *http.Request) {
		session, _ := gws.GetSession(writer, request)

		session.Values["user"] = &UserInfo{
			UserName: "Leon Ding",
			Email:    "ding@ibyte.me",
			Age:      21,
		}

		session.Sync()

		fmt.Fprintln(writer, "set value successful.")
	})

	http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		session, _ := gws.GetSession(writer, request)

		if bytes, ok := session.Values["user"]; ok {
			jsonstr, _ := json.Marshal(bytes)
			fmt.Fprintln(writer, string(jsonstr))
			return
		}

		fmt.Fprintln(writer, "no data")
	})

	http.HandleFunc("/race", func(writer http.ResponseWriter, request *http.Request) {
		session, _ := gws.GetSession(writer, request)
		session.Values["count"] = 0
		var (
			wg  sync.WaitGroup
			mux sync.Mutex
		)
		size := 10000
		wg.Add(size)
		for i := 0; i < size/2; i++ {
			go func() {
				time.Sleep(5 * time.Second)
				mux.Lock()
				if v, ok := session.Values["count"].(int); ok {
					session.Values["count"] = v + 1
				}
				wg.Done()
				mux.Unlock()
			}()
			go func() {
				time.Sleep(5 * time.Second)
				mux.Lock()
				if v, ok := session.Values["count"].(int); ok {
					session.Values["count"] = v + 1
				}
				wg.Done()
				mux.Unlock()
			}()
		}
		wg.Wait()
		fmt.Fprintln(writer, session.Values["count"].(int))
	})

	http.HandleFunc("/result", func(writer http.ResponseWriter, request *http.Request) {
		session, _ := gws.GetSession(writer, request)
		fmt.Fprintln(writer, session.Values["count"].(int))
	})

	_ = http.ListenAndServe(":8080", nil)
}
