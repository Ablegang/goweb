package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"goweb/pkg/logs/loghooks"
	"net/http"
)

func init() {
	// 直接使用 log 包的函数，包内部会实例化一个 std 作为通用 logger
	// std 是指针类型，且是包变量，程序运行时就已经注册

	// 设置最小 log 级别
	log.SetLevel(log.TraceLevel)

	// 设置取堆栈信息
	log.SetReportCaller(true)

	// 注册 Hooks...
	log.AddHook(loghooks.NewEmailNotify())
	log.AddHook(loghooks.NewFileWriter())
}

func main() {
	webServe()
}

func webServe() {
	http.HandleFunc("/", index)              // 设置路由
	err := http.ListenAndServe(":9000", nil) // 监听端口
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// index 路由
func index(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm() // 解析参数，默认不会解析
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println(k, ":", v)
	}
	_, _ = fmt.Fprintf(w, "hello,object!") // 这个写入到 w 的是输出到客户端的
}
