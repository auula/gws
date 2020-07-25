// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/7/25 - 6:56 PM

package main

import (
	"fmt"
	"github.com/higker/go-session"
	"net/http"
)

var manager *session.Manager

func init() {
	manager = session.New(session.MemoryType, "SID", 3000*6000)
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", setHandler)
	http.ListenAndServe(":8080", nil)
}

func IndexHandler(writer http.ResponseWriter, r *http.Request) {
	// init func
	manager.BeginSession(writer, r)
	_, _ = writer.Write([]byte("init session successful!"))
}

func setHandler(writer http.ResponseWriter, r *http.Request) {
	session := manager.GetSessionById(manager.CookieName)
	session.Set("Url", "https://github.com/higker/go-session/")
	_, _ = writer.Write([]byte("set session data successful!"))
}

func getHandler(writer http.ResponseWriter, r *http.Request) {
	session := manager.GetSessionById(manager.CookieName)
	Url := session.Get("Url")
	_, _ = writer.Write([]byte(Url.(string)))
	fmt.Println(Url.(string))
}
