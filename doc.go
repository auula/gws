/*

go—session for Golang session library.
Copyright (c) 2020 SDing <deen.job@qq.com>
Open Source License: MIT License.

	package main

	import (
		"fmt"
		"github.com/higker/go-session"
		"log"
		"net/http"
	)

	func init() {
		// Initializes a Session store that currently supports memory storage
		// that will support Redis or Database in the future
		// cookie参数 --> session.Config{...}
		// CookieName string // sessionID的cookie键名
		// Domain     string // sessionID的cookie作用域名
		// Path       string // sessionID的cookie作用路径
		// MaxAge     int    // 最大生命周期（秒）
		// HttpOnly   bool   // 仅用于http（无法被js读取）
		// Secure     bool   // 启用https

		err := session.Builder(session.MemoryStore, session.DefaultCfg())
		if err != nil {
			log.Fatal(err)
		}
	}

	type user struct {
		Name     string
		password string
	}

	func main() {
		http.HandleFunc("/set", func(writer http.ResponseWriter, request *http.Request) {
			action, err := session.Action(writer, request)
			if err != nil {
				log.Fatal(err)
			}
			action.Set("user", &user{Name: "YNN", password: "password"})
		})

		http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
			action, err := session.Action(writer, request)
			if err != nil {
				log.Fatal(err)
			}
			_, _ = fmt.Fprintln(writer, action.Get("user"))
		})

		http.HandleFunc("/del", func(writer http.ResponseWriter, request *http.Request) {
			action, err := session.Action(writer, request)
			if err != nil {
				log.Fatal(err)
			}
			action.Remove("user")
		})

		http.HandleFunc("/info", func(writer http.ResponseWriter, request *http.Request) {
			action, err := session.Action(writer, request)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(action.ID())
			fmt.Println("======EXEC CLEAR START=====")
			action.Clear()
			fmt.Println("======EXEC CLEAR END=====")
			fmt.Println(action.Get("user"))
		})

		_ = http.ListenAndServe(":6995", nil)
	}

*/

package session
