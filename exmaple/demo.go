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
	err := session.Builder(session.Memory, session.DefaultCfg())
	if err != nil {
		log.Fatal(err)
	}
}

type User struct {
	name string
	age  int8
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
	capture, err := session.Capture(writer, request)
	if err != nil {
		log.Println(err)
	}
	user := User{name: "Ding", age: 21}
	capture.Set("K1", user)
	log.Println("SET User info  OK", user)
	writer.Write([]byte("OK"))
}

func get(writer http.ResponseWriter, request *http.Request) {
	capture, err := session.Capture(writer, request)
	if err != nil {
		log.Println(err)
	}
	bytes, err := capture.Get("K1")
	var u User
	//Deserialize data into objects
	session.DeSerialize(bytes, u)
	log.Println("GET K1 = ", u)
	fmt.Fprintln(writer, u.name, u.age)
}

func clean(writer http.ResponseWriter, request *http.Request) {

}

func del(writer http.ResponseWriter, request *http.Request) {

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
