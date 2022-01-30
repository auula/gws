package main

import (
	"fmt"
	"net/http"

	"github.com/auula/gws"
)

func init() {
	gws.Debug()
	gws.Open(gws.DefaultRAMOption)
}

func main() {
	http.HandleFunc("/set", func(writer http.ResponseWriter, request *http.Request) {

		session, _ := gws.GetSession(writer, request)
		session.Values["foo"] = "bar"
		session.Sync()

		fmt.Fprintln(writer, "set value successful.")
	})

	http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		session, _ := gws.GetSession(writer, request)
		fmt.Fprintln(writer, fmt.Sprintf("foo value is : %s,session create time: %v", session.Values["foo"], session.CreateTime))
	})

	http.HandleFunc("/migrate", func(writer http.ResponseWriter, request *http.Request) {

	})
	_ = http.ListenAndServe(":8080", nil)
}
