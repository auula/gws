# Go-Session

**Session library for Golang🔥.**

## Features

- [x] Session CRUD
- [x] custom config
- [x] simple use

## Use Example

1. go get package

 `go get -u github.com/higker/go-session`
 
2. Example Code

```go
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
			context, err := session.Context(writer, request)
			if err != nil {
				log.Fatal(err)
			}
			context.Set("user", &user{Name: "YNN", password: "password"})
		})
		http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
			context, err := session.Context(writer, request)
			if err != nil {
				log.Fatal(err)
			}
			_, _ = fmt.Fprintln(writer, context.Get("user"))
		})
		http.HandleFunc("/del", func(writer http.ResponseWriter, request *http.Request) {
			context, err := session.Context(writer, request)
			if err != nil {
				log.Fatal(err)
			}
			context.Remove("user")
		})
		http.HandleFunc("/info", func(writer http.ResponseWriter, request *http.Request) {
			context, err := session.Context(writer, request)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(context.ID())
			fmt.Println("======EXEC CLEAR START=====")
			context.Clear()
			fmt.Println("======EXEC CLEAR END=====")
			fmt.Println(context.Get("user"))
		})
		_ = http.ListenAndServe(":6995", nil)
	}
 ```
 3. browser Testing ~
