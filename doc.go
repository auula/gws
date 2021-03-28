// MIT License

// Copyright (c) 2021 Jarvib Ding

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

// Author: SDing <deen.job@qq.com>
// Date: 2021/3/28 - 11:33 PM - UTC/GMT+08:00

// Sessionx is a web session middleware usage example as follows:

// go get github.com/higker/sessionx

// Example Code:

/*
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	sessionx "github.com/higker/sesssionx"
)

var (
	cfg = &sessionx.Configs{
		TimeOut:        time.Minute * 30,
		RedisAddr:      "127.0.0.1:6379",
		RedisDB:        0,
		RedisPassword:  "redis.nosql",
		RedisKeyPrefix: sessionx.SessionKey,
		PoolSize:       100,
		Domain:         "localhost:8080",
		Name:           sessionx.SessionKey,
		Path:           "/",
		Secure:         true,
		HttpOnly:       true,
	}
)

func main() {
	sessionx.New(sessionx.R, cfg)
	http.HandleFunc("/set", func(writer http.ResponseWriter, request *http.Request) {
		session := sessionx.Handler(writer, request)
		session.Set("K", time.Now().Format("2006 01-02 15:04:05"))
		fmt.Fprintln(writer, "set time value succeed.")
	})

	http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		session := sessionx.Handler(writer, request)
		v, err := session.Get("K")
		if err != nil {
			fmt.Fprintln(writer, err.Error())
			return
		}
		fmt.Fprintln(writer, fmt.Sprintf("The stored value is : %s", v))
	})

	http.HandleFunc("/migrate", func(writer http.ResponseWriter, request *http.Request) {
		session := sessionx.Handler(writer, request)
		err := session.MigrateSession()
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintln(writer, session)
	})
	_ = http.ListenAndServe(":8080", nil)
}
*/

package sessionx
