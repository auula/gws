# GWS
**Go's web session library.**

---
[![Go](https://github.com/auula/gws/actions/workflows/go-test.yml/badge.svg?event=push)](https://github.com/auula/gws/actions/workflows/go-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/auula/gws)](https://goreportcard.com/report/github.com/auula/gws)
[![Release](https://img.shields.io/github/v/release/auula/gws.svg?style=flat-square)](https://github.com/auula/gws)
[![License](https://img.shields.io/badge/license-MIT-db5149.svg)](https://github.com/auula/gws/blob/master/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/auula/gws.svg)](https://pkg.go.dev/github.com/auula/gws)
[![codecov](https://codecov.io/gh/auula/gws/branch/dev/graph/badge.svg?token=btbed5BUUZ)](https://codecov.io/gh/auula/gws)
[![DeepSource](https://deepsource.io/gh/auula/gws.svg/?label=active+issues&show_trend=true)](https://deepsource.io/gh/auula/gws/?ref=repository-badge)
[![DeepSource](https://deepsource.io/gh/auula/gws.svg/?label=resolved+issues&show_trend=true)](https://deepsource.io/gh/auula/gws/?ref=repository-badge)


---
[简体中文](./README.md) | [English](./README_EN.md)
---

### Introduction

`GWS` is a `Go` language implementation of a `WEB` `Session` library that supports local session storage, as well as `Redis` remote server distributed storage, and for scalable storage implementations, factories are reserved for developers to implement custom storage to hold session data

### Install

Developers you just need to install this library into your project by executing the following command from within your project.

```shell
go get -u github.com/auula/gws
```

### Use example
You can customize the [`gws.Storage`](./storage.go) interface to use your custom storage, the interface code is as follows:

```go
// Storage global session data store interface.
// You can customize the storage medium by implementing this interface.
type Storage interface {
	// Read data from storage
	Read(s *Session) (err error)
	// Write data to storage
	Write(s *Session) (err error)
	// Remove data from storage
	Remove(s *Session) (err error)
}
```
You only need to implement the `gws.Storage` interface to customize the storage session data, and then register the custom storage implementation interface configuration through the `gws.StoreFactory(opt Options, store Storage)` factory, such as the following example in my demo code inside an example:

```go
package main

import (
	"fmt"
	"net/http"
	// import gws mod
	"github.com/auula/gws"
)

func init() {
	// Whether debug debugging mode is enabled, if so the developer can see the session link log in the console
	// A good developer should look at the logs to analyse the running state of the application, not the debug function in the IDE
	gws.Debug(false)
	// By default configuration and registration of custom storage implementations
	gws.StoreFactory(gws.NewOptions(), &FileStore{})
}

// Customised file storage implementations
type FileStore struct{}

func (fs FileStore) Read(s *gws.Session) (err error) {
	panic("implement me")
}

func (fs FileStore) Write(s *gws.Session) (err error) {
	panic("implement me")
}

func (fs FileStore) Remove(s *gws.Session) (err error) {
	panic("implement me")
}

func main() {
	// testing file store
	http.HandleFunc("/panic", func(writer http.ResponseWriter, request *http.Request) {
		// gws.GetSession will return the session
		session, _ := gws.GetSession(writer, request)
		// Save the data that needs to be stored for the session via session.Values
		session.Values["foo"] = "bar"
		// The Sync method is used to synchronise the data persistence, although it can be left out if it is the default in-memory storage.
		// If it is a remote server or custom storage be sure to perform this method to synchronize data to other distributed ends
		session.Sync()

		fmt.Fprintln(writer, "set value successful.")
	})
	_ = http.ListenAndServe(":8080", nil)

}
```
The above example code, showing how to customise a storage implementation, can be found at: [example/store_example.go](./example/store_example.go)

---
If you are using a single machine, or a small `Web Service` application, you can use the default local `in-memory` storage, where the session storage is stored in the local server `memory`, the downside of this is that the session data cannot be recovered when the application is restarted.

If you want to support persistence you can customise it, or you can use the `Redis` solution provided by gws by default to solve the session distributed storage, `Redis` in a single point of failure, as long as the session has not expired, the data will still be synchronised after the application node is restored to normal.

```go
func init() {
	// Customize the configuration options, see the documentation on go.dev for details, or look at the source code
	// var opt gws.RAMOption
	// opt.Domain = "www.ibyte.me"
	// gws.Open(opt)

	// You can use the default configuration initialization, initialized by option function mode
	
	// gws.Open(gws.NewOptions())
	// gws.Open(gws.NewOptions(gws.Domain(""), gws.CookieName("")))

	// Recommended direct default configuration
	gws.Open(gws.DefaultRAMOptions)

	// This is to initialize the Redis distributed storage
	// gws.Open(gws.NewRDSOptions("127.0.0.1", 6379, "redis.nosql"))
}
```
In the following sample code, I will demonstrate how to manage your session data via `gws`:

```go
// Example Data
type UserInfo struct {
	UserName string `json:"user_name,omitempty"`
	Email    string `json:"email,omitempty"`
	Age      uint8  `json:"age,omitempty"`
}
```

I configured a  `/set` route, how to store a user value inside the session, store the value directly using the `Values` field assignment, in fact it is a `map[string]interface{}` variant structure, note that the `Values` here is not parallel safe, in fact I considered this in the development of gws and wanted to design concurrency safe In fact, I had this in mind when developing gws and wanted to design concurrency-safe api's, but considered that too many api's would not be good, and that writing `Go` should be kept simple, not like `Java` where you have to abstract through `get` and `set`, which would just make your code base huge and cluttered.

So in the documentation I make it clear that if you are concurrently manipulating Values and custom locking! Example code will also be added later.

```go
http.HandleFunc("/set", func(writer http.ResponseWriter, request *http.Request) {

	session, _ := gws.GetSession(writer, request)
	session.Values["user"] = &UserInfo{
		UserName: "Leon Ding",
		Email:    "ding@ibyte.me",
		Age:      21,
	}

	// The ram mode can be executed without it, as it is a memory pointer reference
	session.Sync()

	fmt.Fprintln(writer, "set value successful.")
})

```
If you want to read data from the session, you can see the sample code:

```go
http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
	session, _ := gws.GetSession(writer, request)

	// Reading data is the same operation as detecting map, you can omit this if operation if you can be sure that this must have a value
	if bytes, ok := session.Values["user"]; ok {
		jsonstr, _ := json.Marshal(bytes)
		fmt.Fprintln(writer, string(jsonstr))
		return
	}

	fmt.Fprintln(writer, "no data")
})
```
Delete operation:

```go
http.HandleFunc("/del", func(writer http.ResponseWriter, request *http.Request) {
	session, _ := gws.GetSession(writer, request)
	delete(session.Values, "user")
	// Must be synchronized, if it's a custom store or a Redis distributed store
	session.Sync()
	fmt.Fprintln(writer, "successful")
})
```

Now that you know the basics of how to use it, you can easily manage your session data. You can view the sample code to do your job better, or view the source code.

## Coroutine Safe

As I designed the API without the intention of writing a `get、 set、 del` and then providing an internal lock in it to ensure and secure that all callers must lock themselves in case of data competition, or you go write the slip that you can customize to be wrapped and secure, the following code I demonstrate how to concurrently secure the operation.


```go
http.HandleFunc("/race", func(writer http.ResponseWriter, request *http.Request) {
	session, _ := gws.GetSession(writer, request)

	session.Values["count"] = 0
	var (
		wg  sync.WaitGroup
		// coroutine safe lock
		mux sync.Mutex
	)
	size := 10000
	wg.Add(size)
	for i := 0; i < size/2; i++ {
		go func() {
			time.Sleep(5 * time.Second)
			mux.Lock()
			if v, ok := session.Values["count"].(int); ok {
				session.Values["count"] = v + 1
			}
			wg.Done()
			mux.Unlock()
		}()
		go func() {
			time.Sleep(5 * time.Second)
			mux.Lock()
			if v, ok := session.Values["count"].(int); ok {
				session.Values["count"] = v + 1
			}
			wg.Done()
			mux.Unlock()
		}()
	}
	wg.Wait()
	fmt.Fprintln(writer, session.Values["count"].(int))
})
```

In a data contention state, other callers can fetch the value normally, but you have to ensure that you customize the lock control range, what type of lock you want to use, such as `read/write lock` or mutually exclusive lock, this depends on your understanding of go, or you are very strong, you can solve the data contention through the channel, when I design the API, I am keeping as little as possible to affect or limit the caller some operation experience, the above and the following examples are in the case of race in the request, the result will not block and can still fetch the value of the demonstration:

```go
http.HandleFunc("/result", func(writer http.ResponseWriter, request *http.Request) {
	session, _ := gws.GetSession(writer, request)
	fmt.Fprintln(writer, session.Values["count"].(int))
})
```

## Session Hijacking

Session fixed attack, this process, normal users in the browser to access the site we write, but this time there is a `hack` through `arp` spoofing, the router traffic hijacked to his computer, and then the hacker through some special software to grab your network request traffic information, in this process if you `sessionid` if stored in the cookie, it is likely to be If you log in to the site at this time, this is the hacker to get your login credentials, and then in the login to replay that is to use your `sessionid`, so as to achieve the purpose of access to your account-related data.


To do this I added a `gws.Migrate(write http.ResponseWriter, old *Session) (*Session, error) `built-in function to gws, using the following example：


```go
http.HandleFunc("/migrate", func(writer http.ResponseWriter, request *http.Request) {
	var (
		session *gws.Session
		err     error
	)

	session, _ = gws.GetSession(writer, request)
	log.Printf("old session %p \n", session)

	// Migrate session data and refresh client sessions, discarding old sessions
	if session, err = gws.Migrate(writer, session); err != nil {
		fmt.Fprintln(writer, err.Error())
		return
	}

	log.Printf("old session %p \n", session)
	jsonstr, _ := json.Marshal(session.Values["user"])
	fmt.Fprintln(writer, string(jsonstr))
})
```

`Migrate` will help you migrate session data, you can also use it with the https protocol, but of course the API I had in mind when I designed `gws`, all of it is provided.

#### Table of contents for the above example code：

- [./example/store_example.go](./example/store_example.go)
- [./example/ram_example.go](./example/ram_example.go)
- [./example/redis_example.go](./example/redis_example.go)

