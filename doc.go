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

// Go's web session library.
// gws is a go language implementation of the web session library,
// which supports local storage of session data as well as redis
// remote distributed storage, distributed session sharing.

// author: Leon Ding <ding@ibyte.me>

// Quick Start
// 1. install mod
// go get github.com/auula/gws
// 2. Use Example Code

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/auula/gws"
// )

// func init() {

// 	gws.Open(gws.DefaultRAMOptions)

// }

// type UserInfo struct {
// 	UserName string `json:"user_name,omitempty"`
// 	Email    string `json:"email,omitempty"`
// 	Age      uint8  `json:"age,omitempty"`
// }

// func main() {
// 	http.HandleFunc("/set", func(writer http.ResponseWriter, request *http.Request) {
// 		session, _ := gws.GetSession(writer, request)

// 		session.Values["user"] = &UserInfo{
// 			UserName: "Leon Ding",
// 			Email:    "ding@ibyte.me",
// 			Age:      21,
// 		}
// 		session.Sync()

// 		fmt.Fprintln(writer, "set value successful.")
// 	})

// 	http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
// 		session, _ := gws.GetSession(writer, request)

// 		if bytes, ok := session.Values["user"]; ok {
// 			jsonstr, _ := json.Marshal(bytes)
// 			fmt.Fprintln(writer, string(jsonstr))
// 			return
// 		}

// 		fmt.Fprintln(writer, "no data")
// 	})

// 	http.HandleFunc("/userinfo", func(writer http.ResponseWriter, request *http.Request) {
// 		session, err := gws.GetSession(writer, request)
// 		if err != nil {
// 			fmt.Fprintln(writer, err.Error())
// 			return
// 		}
// 		jsonstr, _ := json.Marshal(session.Values["user"])
// 		fmt.Fprintln(writer, string(jsonstr))
// 	})

// 	http.HandleFunc("/migrate", func(writer http.ResponseWriter, request *http.Request) {
// 		var (
// 			session *gws.Session
// 			err     error
// 		)

// 		session, _ = gws.GetSession(writer, request)
// 		log.Printf("old session %p \n", session)

// 		if session, err = gws.Migrate(writer, session); err != nil {
// 			fmt.Fprintln(writer, err.Error())
// 			return
// 		}

// 		log.Printf("new session %p \n", session)
// 		jsonstr, _ := json.Marshal(session.Values["user"])
// 		fmt.Fprintln(writer, string(jsonstr))
// 	})

// 	_ = http.ListenAndServe(":8080", nil)
// }

package gws
