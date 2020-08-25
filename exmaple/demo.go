// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/24 - 7:39 下午 - UTC/GMT+08:00

package main

import (
	"fmt"
	"github.com/higker/go-session"
	"log"
	"net/http"
)

func init() {
	//cfg := session.Config{
	//	CookieName:     session.DefaultCookieName,
	//	Path:           "/",
	//	MaxAge:         60,
	//	HttpOnly:       true,
	//	Secure:         false,
	//	RedisAddr:      "128.199.155.162:6379",
	//	RedisPassword:  "deen.job",
	//	RedisDB:        0,
	//	RedisKeyPrefix: session.RedisPrefix,
	//}
	//err := session.Builder(session.Redis, &cfg)
	//if err != nil {
	//	log.Fatal(err)
	//}
	err := session.Builder(session.Memory, session.DefaultCfg())
	if err != nil {
		log.Fatal(err)
	}
}

type User struct {
	Name string `json:"name"`
	Age  int8   `json:"age"`
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/set", set)
	http.HandleFunc("/get", get)
	http.HandleFunc("/del", del)
	http.HandleFunc("/clean", clean)
	http.ListenAndServe(":8080", nil)
}

func set(writer http.ResponseWriter, request *http.Request) {
	ctx, err := session.Ctx(writer, request)
	if err != nil {
		log.Println(err)
	}
	// set data for session
	user := User{Name: "Ding", Age: 21}
	ctx.Set("K1", user)
	fmt.Fprintln(writer, "set value ok")
}

func get(writer http.ResponseWriter, request *http.Request) {
	ctx, err := session.Ctx(writer, request)
	if err != nil {
		log.Println(err)
	}
	bytes, err := ctx.Get("K1")
	if err != nil {
		log.Println("ERR", err)
	}
	u := new(User)
	fmt.Println(bytes)
	//Deserialize data into objects
	session.DeSerialize(bytes, u)
	log.Println("GET K1 = ", u)
	fmt.Fprintln(writer, u)
}

func clean(writer http.ResponseWriter, request *http.Request) {
	ctx, err := session.Ctx(writer, request)
	if err != nil {
		log.Println(err)
	}
	// clean session data
	ctx.Clean()
	fmt.Fprintln(writer, "clean data ok")
}

func del(writer http.ResponseWriter, request *http.Request) {
	ctx, err := session.Ctx(writer, request)
	if err != nil {
		log.Println(err)
	}
	err = ctx.Del("K1")
	if err != nil {
		fmt.Fprintln(writer, err)
	}
	fmt.Fprintln(writer, "delete v1 successful")
}

func index(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(writer, `
        Go session storage example:<br><br>
        <a href="/set">Store key in session</a><br>
        <a href="/get">Get key value from session</a><br>
        <a href="/del">Destroy session</a>
		<a href="/clean">Clean session</a>
		<a href="https://github.com/higker/go-session">to github</a><br>`)
}
