# go-session
- Session library for go language.
- 这是一个开源的`Go`语言第三方库
- 由于`Go`官方没有提供原生的`Session`库,所以笔者写了这么一个库，希望你在开发时候的使用这个库,减少工作量,少写几行代码，😁。
- By: SDing 2020-07-30 20:17:32

# Get this Package

`go get -u github.com/higker/go-session`

# Use Example Code

```go
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
	handel, err := session.Handle(writer, r)
	if err != nil {
		fmt.Println(err)
	}
	// SET data
	handel.Set("Key","https://github.com/higker/go-session/")
	_, _ = writer.Write([]byte("init session successful!"))
}


func getHandler(writer http.ResponseWriter, r *http.Request) {
	handel, err := session.Handle(writer, r)
	if err != nil {
		fmt.Println(err)
	}
	// GET data
	_, _ = fmt.Fprintln(writer, handel.Get("Key"))
	fmt.Println(handel.GetID())
}
```