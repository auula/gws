# go-session
 ![Go](https://github.com/airplayx/gormat/workflows/Go/badge.svg)
 [![Go Report Card](https://goreportcard.com/badge/github.com/airplayx/gormat)](https://goreportcard.com/report/github.com/higker/go-session)
 [![shields](https://img.shields.io/github/v/release/higker/go-session.svg)](https://github.com/higker/go-session/releases)
 
**Session library for Golang🔥.**


## Features

- [x] Simple Use
- [x] Session CRUD
- [x] Custom Config
- [x] Memory Storage
- [x] Redis Storage
- [x] Distributed session

## Warning

- 目前本库已经实现内存存储和Redis做分布式存储
- 内存版本使用的是计算机内存 所以可能在运行的时候内存可能会大一点
- session从第一次请求过来就创建 生命周期根据你自己设置单位秒
-  每次都是按照生命周期来计算一个session的周期
- 例如一次30分钟那这个session周期就是30分钟
- 中途你做了set或者get也没有用，到了时间这个session就自动销毁了，和其他session库有一点不一样！！！！

## Use Example

> Go version => 1.11

1.Get Package

 `go get -u github.com/higker/go-session`
 
2.Example Code

```go
package main

import (
	"fmt"
	"github.com/higker/go-session"
	"log"
	"net/http"
)

func init() {
	cfg := session.Config{
		CookieName:     session.DefaultCookieName,
		Path:           "/",
		MaxAge:         session.DefaultMaxAge,
		HttpOnly:       true,
		Secure:         false,
		RedisAddr:      "128.199.155.162:6379",
		RedisPassword:  "your password",
		RedisDB:        0,
		RedisKeyPrefix: session.RedisPrefix,
	}
	err := session.Builder(session.Redis, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Memory store 使用内存存储方式就使用下面这个 注释上面的 打开下面的
	// It is not recommended that you use it because it consumes memory
	//err := session.Builder(session.Memory, session.DefaultCfg())
	//if err != nil {
	//	log.Fatal(err)
	//}
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
	//Deserialize data into objects
	session.DeSerialize(bytes, u)
	fmt.Fprintln(writer, u)
}

func clean(writer http.ResponseWriter, request *http.Request) {
	ctx, err := session.Ctx(writer, request)
	if err != nil {
		log.Println(err)
	}

	// clean session all data by session
	ctx.Clean(writer)

	fmt.Fprintln(writer, "clean data ok")
}

func del(writer http.ResponseWriter, request *http.Request) {
	ctx, err := session.Ctx(writer, request)
	if err != nil {
		log.Println(err)
	}
	err = ctx.Del("K1")
	if err != nil {
		log.Println(err)
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

 ```
3.browser Testing ~  Good luck😜~
 > 由于go官方没有提供session的标准库，所以笔者自己写了一个并且开源出来希望你帮助屏幕前需要的你，给个star吧~
