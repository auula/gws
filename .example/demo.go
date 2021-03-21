package main

import (
	"fmt"
	sessionx "github.com/higker/sesssionx"
	"net/http"
	"time"
)

var (
	cfg = &sessionx.Configs{
		EncryptedKey:    "0123456789012345",
		SessionLifeTime: time.Minute * 30,
		Cookie: &http.Cookie{
			Name:     sessionx.SESSION_KEY,
			Path:     "/",
			Expires:  time.Now().Add(time.Minute * 30),
			Secure:   false,
			HttpOnly: true,
			MaxAge:   60 * 30,
		},
	}
)

func main() {
	sessionx.New(sessionx.M, cfg)
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
