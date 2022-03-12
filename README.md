# GWS

**Go's web session library.**

---
[![Go](https://github.com/auula/gws/actions/workflows/go-test.yml/badge.svg?event=push)](https://github.com/auula/gws/actions/workflows/go-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/auula/gws)](https://goreportcard.com/report/github.com/auula/gws)
[![Release](https://img.shields.io/github/v/release/auula/gws.svg?style=flat-square)](https://github.com/auula/gws)
[![License](https://img.shields.io/badge/license-MIT-db5149.svg)](https://github.com/auula/gws/blob/master/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/auula/gws.svg)](https://pkg.go.dev/github.com/auula/gws)
[![codecov](https://codecov.io/gh/auula/gws/branch/dev/graph/badge.svg?token=btbed5BUUZ)](https://codecov.io/gh/auula/gws)
[![DeepSource](https://deepsource.io/gh/auula/gws.svg/?label=active+issues&token=c0sepC4oiYxaeLOxgJFLfSWP)](https://deepsource.io/gh/auula/gws/?ref=repository-badge)

---


[简体中文](./README.md) | [English](./README_EN.md)

---
### 介 绍

`GWS`是一个`Go`语言实现的`WEB`会话库，支持本地会话存储，也支持`Redis`远程服务器分布式存储，并且为了可扩展存储实现，预留工厂，方便开发者自定义实现存储来保存会话数据。

### 安 装
开发者你只需要安装本库到你到项目里面，在你的项目里面执行下面命令即可安装：
```shell
go get -u github.com/auula/gws
```
### 使用示例

首先要声明一点，`gws`是支持多种存储介质保存`session`数据的，你可以自定义实现[`gws.Storage`](./storage.go)存储接口，来使用你自定义存储，接口代码如下:

```go
// Storage global session data store interface.
// You can customize the storage medium by implementing this interface.
type Storage interface {
	// Read data from store
	Read(s *Session) (err error)
	// Write data to storage
	Write(s *Session) (err error)
	// Remove data from storage
	Remove(s *Session) (err error)
}
```
你只需要实现[`gws.Storage`](./storage.go)接口，就可以自定义存储会话数据，然后通过`gws.StoreFactory(opt Options, store Storage)`工厂注册自定义存储实现接口配置，例如下面我在演示代码里面一个例子:

```go
package main

import (
	"fmt"
	"net/http"
	// 导入gws模块
	"github.com/auula/gws"
)

func init() {
	// 是否开启debug调试模式，如果开启则开发者可以在控制台看到会话链路日志
	// 好的开发者应该看日志去分析程序运行状态，而不是集成开发环境里面的debug功能
	gws.Debug(false)
	// 通过默认配置，并且注册自定义存储实现
	gws.StoreFactory(gws.NewOptions(), &FileStore{})
}

// 自定义的文件存储实现
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
	// 测试自定义存储
	http.HandleFunc("/panic", func(writer http.ResponseWriter, request *http.Request) {
		// gws.GetSession 会返回本次请求的session
		session, _ := gws.GetSession(writer, request)
		// 通过session.Values 保存需要存储会话的数据
		session.Values["foo"] = "bar"
		// 通过Sync方法同步数据持久化，当然这里如果是默认内存存储可以不调用
		// 如果是远程服务器或者自定义存储一定要执行此方法同步数据到其他分布式端
		session.Sync()

		fmt.Fprintln(writer, "set value successful.")
	})

	_ = http.ListenAndServe(":8080", nil)

}
```
以上示例代码，展示如何自定义实现一个存储，具体示例代码请查看：[./example/store_example.go](./example/store_example.go)

---
如果只是单机使用，或者是一个小体积`Web Service`应用，你可以使用默认的本地内存存储，会话存储会保存在本地服务器内存里面，这个缺点就是程序重启会话数据不能恢复。

如果想支持持久化你可以自定义，也可以使用`gws`默认提供的`Redis`方案去解决会话分布式存储，`Redis`在单点故障时，只要会话没有过期，应用节点恢复正常之后，数据依旧会同步到。

```go
func init() {
	// 自定义配置选项参数，具体哪些参数可以查看go.dev上面的文档，或者看源代码吧
	// var opt gws.RAMOption
	// opt.Domain = "www.ibyte.me"
	// gws.Open(opt)

	// 你可以使用默认配置初始化，通过option function模式初始化
	
	// gws.Open(gws.NewOptions())
	// gws.Open(gws.NewOptions(gws.Domain(""), gws.CookieName("")))

	// 推荐直接默认配置
	gws.Open(gws.DefaultRAMOptions)

	// 这个是初始化Redis分布式存储的
	// gws.Open(gws.NewRDSOptions("127.0.0.1", 6379, "redis.nosql"))
}
```

下面的示例代码，我会演示如何通过`gws`管理你的会话数据：

```go
// 为了演示数据变化，我定义的一个UserInfo结构体
type UserInfo struct {
	UserName string `json:"user_name,omitempty"`
	Email    string `json:"email,omitempty"`
	Age      uint8  `json:"age,omitempty"`
}
```
我配置了一个`set`路由，如何在会话里面存储一个`user`的值，存储值直接使用`Values`字段赋值，其实就是一个`map[string]interface{}`变体结构，注意这里的`Values`不是并行安全的，其实我在开发`gws`就考虑到了这个问题，并且想设计并发安全的`api`，但是考虑到`api`太多了也不好，写`Go`要保持大道至简，并不是像`Java`那样要通过`get`和`set`各种抽象，那样只会让你的代码库变得庞大，杂乱无章。

所以在文档我明确说明了如果是并发操作`Values`并且自定义加锁！！！示例代码也会在后面添加：
```go
http.HandleFunc("/set", func(writer http.ResponseWriter, request *http.Request) {

	session, _ := gws.GetSession(writer, request)
	session.Values["user"] = &UserInfo{
		UserName: "Leon Ding",
		Email:    "ding@ibyte.me",
		Age:      21,
	}

	// ram模式可以不用执行，因为是内存指针引用
	session.Sync()

	fmt.Fprintln(writer, "set value successful.")
})
```
如果要从会话里面读取数据，可以看示例代码：

```go
http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
	session, _ := gws.GetSession(writer, request)

	// 读取数据和检测map一样的操作，你如果能确保这个一定有值，你可以省去这个if操作
	// 直接取值也行
	if bytes, ok := session.Values["user"]; ok {
		jsonstr, _ := json.Marshal(bytes)
		fmt.Fprintln(writer, string(jsonstr))
		return
	}

	fmt.Fprintln(writer, "no data")
})
```
删除操作及其简单，如果你是老司机开发者，我相信你已经不需要看示例代码了，如下：

```go
http.HandleFunc("/del", func(writer http.ResponseWriter, request *http.Request) {
	session, _ := gws.GetSession(writer, request)
	delete(session.Values, "user")
	// 一定同步，如果是自定义存储或者是Redis分布式存储的话
	session.Sync()
	fmt.Fprintln(writer, "successful")
})
```

如果要清理这个`session`的`Values`数据，可以使用`gws.Malloc(v *gws.Values)`函数：

```go
http.HandleFunc("/clean", func(rw http.ResponseWriter, request *http.Request) {
	session, _ := gws.GetSession(rw, request)
	// clean session data
	gws.Malloc(&session.Values)
	// sync session modify
	session.Sync()
	fmt.Fprintf(rw, "clean session data successful.")
})
```

如果废弃掉这个`session`则调用`gws.Invalidate(s *Session) error`函数：

```go
http.HandleFunc("/invalidate", func(rw http.ResponseWriter, request *http.Request) {
	session, _ := gws.GetSession(rw, request)
	gws.Invalidate(session)
	fmt.Fprintf(rw, "set session invalidate successful.")
})
```
上面都是基本的增删改查操作，如果你作为一名`API`调用工程师或者是`API`操作员，那你看到这估计就差不多了，可以完成你日常的开发需求了，你也不需要去了解内部实现，如果要了解内部实现，我后面有空会去讲内部实现。


## 并行安全
由于我在设计`API`的时候，没有打算去写一个`get、set、del`，然后在里面提供一个内部锁去保证并行安全，所有调用者必须在有数据竞争的情况下自行加锁，或者你`go`写的溜，你可以自定义去包装并行安全的，下面这段代码我演示了如何并行安全的操作：
```go
http.HandleFunc("/race", func(writer http.ResponseWriter, request *http.Request) {
	session, _ := gws.GetSession(writer, request)

	session.Values["count"] = 0
	var (
		wg  sync.WaitGroup
		// 如果你是并发操作请加锁
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
在数据竞争状态下，其他的调用者，可以正常取值，但是你要保证你自定义的锁的控制范围，你要用什么类型的锁，例如读写锁还是互斥锁，这个要看你对`go`的了解程度了，或者你很强，你可以通过`channel`解决数据竞争，我在设计`API`的时候，我是保持着尽可能少量的去影响或者限制调用者一些操作体验的，上面和下面的示例是在`race`在请求的情况下，`result`不会阻塞并且还能取值的演示：

```go
http.HandleFunc("/result", func(writer http.ResponseWriter, request *http.Request) {
	session, _ := gws.GetSession(writer, request)
	fmt.Fprintln(writer, session.Values["count"].(int))
})
```


## 会话劫持

会话固定攻击，这个过程，正常用户在通过浏览器访问我们编写的网站，但是这个时候有个`hack`通过`arp`欺骗，把路由器的流量劫持到他的电脑上，然后黑客通过一些特殊的软件抓包你的网络请求流量信息，在这个过程中如果你`sessionid`如果存放在`cookie`中，很有可能被黑客提取处理，如果你这个时候登录了网站，这是黑客就拿到你的登录凭证，然后在登录进行重放也就是使用你的`sessionid`，从而达到访问你账户相关的数据目的。

为此我在`gws`里面添加一个`gws.Migrate(write http.ResponseWriter, old *Session) (*Session, error)`内置函数，使用示例：

```go
http.HandleFunc("/migrate", func(writer http.ResponseWriter, request *http.Request) {
	var (
		session *gws.Session
		err     error
	)

	session, _ = gws.GetSession(writer, request)
	log.Printf("old session %p \n", session)

	// 迁移会话数据，并且刷新客户端会话，丢弃掉老的session
	if session, err = gws.Migrate(writer, session); err != nil {
		fmt.Fprintln(writer, err.Error())
		return
	}

	log.Printf("old session %p \n", session)
	jsonstr, _ := json.Marshal(session.Values["user"])
	fmt.Fprintln(writer, string(jsonstr))
})
```
`gws.Migrate`会帮助你迁移会话数据，你也可以配合`https`协议使用，当然该有的`API`我在设计`gws`的时候就已经考虑到了，所有都提供了。

**以上示例代码目录：**

- [./example/store_example.go](./example/store_example.go)
- [./example/ram_example.go](./example/ram_example.go)
- [./example/redis_example.go](./example/redis_example.go)

如果你发现了什么`bug`欢迎`pr`或者`issues`~如果对你有帮助你可以按一个`star`再走呗，[https://github.com/auula/gws](https://github.com/auula/gws)
