package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	sessionx "github.com/higker/sessionx"
)

var (
	cfg = &sessionx.Configs{
		TimeOut:        time.Minute * 30,
		RedisAddr:      "127.0.0.1:6379",
		RedisDB:        0,
		RedisPassword:  "redis.nosql",
		RedisKeyPrefix: sessionx.SessionKey,
		PoolSize:       100,
		Domain:         "localhost", // set domain by you
		Name:           sessionx.SessionKey,
		Path:           "/",
		Secure:         true,
		HttpOnly:       true,
	}
	// 如果使用内存存储就直接使用 sessionx.DefaultCfg
	// sessionx.New(sessionx.M, sessionx.DefaultCfg)
)

func main() {
	sessionx.New(sessionx.M, sessionx.DefaultCfg)
	http.HandleFunc("/set", func(writer http.ResponseWriter, request *http.Request) {
		session := sessionx.Handler(writer, request)
		// 存储K的值
		session.Set("K", time.Now().Format("2006 01-02 15:04:05"))
		fmt.Fprintln(writer, "set time value succeed.")
	})

	http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		session := sessionx.Handler(writer, request)
		// 获取存储K的值
		v, err := session.Get("K")
		if err != nil {
			fmt.Fprintln(writer, err.Error())
			return
		}
		fmt.Fprintln(writer, fmt.Sprintf("The stored value is : %s", v))
	})

	http.HandleFunc("/migrate", func(writer http.ResponseWriter, request *http.Request) {
		session := sessionx.Handler(writer, request)

		// MigrateSession 函数会迁移session会话数据，返回新的session
		session, err := session.MigrateSession()
		if err != nil {
			log.Println(err)
		}

		session.Set("person", "Jarvib Ding")
		fmt.Fprintln(writer, session)
	})
	_ = http.ListenAndServe(":8080", nil)
}
