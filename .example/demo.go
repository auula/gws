package main

import (
	"fmt"
	sessionx "github.com/higker/sesssionx"
	"net/http"
	"time"
)

var (
	cfg = &sessionx.Configs{
		TimeOut:        time.Minute * 30,
		RedisAddr:      "127.0.0.1:6379",
		RedisDB:        0,
		RedisPassword:  "redis.nosql",
		RedisKeyPrefix: sessionx.SessionKey,
		PoolSize:       100,
		Cookie: &http.Cookie{
			Name:     sessionx.SessionKey,
			Path:     "/",
			Expires:  time.Now().Add(time.Minute * 30),
			Secure:   false,
			HttpOnly: true,
		},
	}
)

func main() {
	sessionx.New(sessionx.R, cfg)
	http.HandleFunc("/set", func(writer http.ResponseWriter, request *http.Request) {
		session := sessionx.Handler(writer, request)
		_ = session.Set("K", time.Now().Format("2006 01-02 15:04:05"))
		_, _ = fmt.Fprintln(writer, "set succeed.")
	})
	http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		session := sessionx.Handler(writer, request)
		v, _ := session.Get("K")
		_, _ = fmt.Fprintln(writer, v)
	})
	_ = http.ListenAndServe(":8080", nil)
}
