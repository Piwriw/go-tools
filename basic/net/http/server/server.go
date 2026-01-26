package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	//http://127.0.0.1:8000/go
	// 单独写回调函数
	http.HandleFunc("/go", myHandler)
	http.HandleFunc("/get", getParamsHandler)
	http.HandleFunc("/post", postHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/Users/joohwan/GolandProjects"))))
	//http.HandleFunc("/ungo",myHandler2 )
	// addr：监听的地址
	// handler：回调函数
	http.ListenAndServe("127.0.0.1:8000", nil)
}

// postHandler Post
func postHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// 1. 请求类型是application/x-www-form-urlencoded时解析form数据
	r.ParseForm()

	fmt.Println(r.PostForm) // 打印form数据
	w.Write([]byte(fmt.Sprintf("name:%s,age:%s", r.PostForm.Get("name"), r.PostForm.Get("age"))))
	fmt.Println(r.PostForm.Get("name"), r.PostForm.Get("age"))

	// 2. 请求类型是application/json时从r.Body读取数据
	b, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("read request.Body failed, err:%v\n", err)
		return
	}

	w.Write(b)
}

// handler函数
func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr, "连接成功")
	// 请求方式：GET POST DELETE PUT UPDATE
	fmt.Println("method:", r.Method)
	// /go
	fmt.Println("url:", r.URL.Path)
	fmt.Println("header:", r.Header)
	fmt.Println("body:", r.Body)
	// 回复
	w.Write([]byte("www.5lmh.com"))
}

func getParamsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := r.URL.Query()
	name := data.Get("name")
	age := data.Get("age")
	res := fmt.Sprintf("name:%s,age:%s", name, age)
	w.Write([]byte(res))
}
