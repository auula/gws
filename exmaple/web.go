// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/7/30 - 8:06 PM

package main

import (
	"fmt"
	"github.com/higker/go-session"
	"net/http"
)

func init() {
	// 初始化一个Session存储器 目前支持内存存储 未来将支持 Redis 或者 Database
	err := session.Builder(session.NewMemoryStore())
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)
	_ = http.ListenAndServe(":8080", nil)
}

func setHandler(writer http.ResponseWriter, r *http.Request) {
	// init func
	handel, err := session.Handel(writer, r)
	if err != nil {
		fmt.Println(err)
	}
	// SET data
	handel.Set("Key", "https://github.com/higker/go-session/")
	_, _ = writer.Write([]byte("init session successful!"))
}

func getHandler(writer http.ResponseWriter, r *http.Request) {
	handel, err := session.Handel(writer, r)
	if err != nil {
		fmt.Println(err)
	}
	// GET data
	_, _ = fmt.Fprintln(writer, handel.Get("Key"))
	fmt.Println(handel.GetID())
}
