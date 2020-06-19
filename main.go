package main

import (
	"fmt"
	"goweb/pkg/logs/loghooks"
	"net/http"
	"os"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	// 直接使用 log 包的函数，包内部会实例化一个 std 作为通用 logger
	// std 是指针类型，且是包变量，程序运行时就已经注册
	// 设置日志格式
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano, // 含纳秒
	})

	log.SetLevel(log.TraceLevel)

	// 注册 Hooks...
	log.AddHook(loghooks.NewEmailNotify())
	log.AddHook(loghooks.NewFileWriter())
	log.AddHook(&loghooks.FileWriter{
		LogMode:          "single",
		Dir:              "storage/logs-test/",
		FileNameFormater: "logs.txt",
		EntryFormatter: &log.JSONFormatter{
			TimestampFormat: time.RFC3339Nano, // 含纳秒
		},
		HookLevels: log.AllLevels,
		Perm:       os.FileMode(0777),
	})
}

func main() {
	log.WithField("xxxxxx", "9999999").Infoln("9999999999")
	log.Debugln("88888")
	log.Traceln("-1-1-1-1-1-1-1-1-")
	log.Exit(2)
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

// GO 基础 ...
func basic() {

	// 特点：
	// 并发、带 GC、快速编译、静态类型（轻量级类型系统）、编译型
	// 类型系统没有层级
	// 源码安装、标准包安装（一般开发环境都这么安装）、第三方工具安装（wget、apt-get）
	// 一个系统多个 go 版本，需要用到三方工具 GVM

	// 关键环境变量：
	// GOROOT go 的安装路径
	// GOPATH 源码与 import 包的搜索路径，一般约定三个目录，src（源码文件）、pkg（编译后生成的文件）、bin（编译后生成的可执行文件）
	// PATH 一般是 GOROOT 和 GOPATH 的合并
	// GOBIN 不允许设置多个路径，可以不设置。为空时遵循约定由于配置原则，可执行文件放置各自GOPATH目录的bin文件夹中
	// 不同操作系统按自己的方式去设置环境变量

	// 目录结构规划：
	// $GOPATH/src 下一般存放了所有项目的源码，子目录按项目名划分，如 $GOPATH/src/blog 、$GOPATH/src/web、$GOPATH/src/app
	// 也支持多级目录，如 $GOPATH/src/github.com/astaxie/beedb ，这样在其它项目中以 "github.com/astaxie/beedb" 来引入包
	// 使用 go mod 的话，使用 go mod 引入的包会 $GOPATH/pkg/mod 下

	// go get
	// 在没有使用 go mod 的项目里，使用 go get 来获取远程包
	// go get 支持 github、googlecode 等远程仓库

	// go build
	// go clean 移除当前源码包和关联源码包里编译生成的文件，一般通过这个命令清理编译文件，再向 github 提交代码
	// go fmt 格式化代码，go 本身强制要求代码格式，可以通过 go fmt 命令来整理代码格式
	// go install
	// go test
	// go tool fix . 修复老版本代码到新版本
	// go tool vet 分析当前目录或文件是否都是正确的代码
	// go version
	// go env
	// go list
	// go run

	// go 关键字：
	// break default func interface select case
	// defer go map struct chan else goto package switch
	// const fallthrough if range type continue for import return var

	// go 天生支持 utf-8，main.main() 是项目的入口
	// := 只能在函数内部使用
	// 大写字母开头为公有，小写字母开头为私有，这个相对包而言
	// [4]int 和 [5]int 是两种类型的数组，长度本身就是数组的一部分，所以数组本身长度是固定的
	// slice 是一个引用类型，总是指向底层的一个 array
	// ar[0:len(ar)] 可以缩写为 ar[:]， 0 和 len 都是可以省略的
	// len、cap、append、copy

	// 在 go 中没有值是可以安全地并发读写，它不是 thread-safe 的，所以在多个 goroutine 进行读写时，一定要用 mutex lock 机制

	// map 有两个返回值
	dt := map[string]string{
		"name": "ben",
		"age":  "18",
	}
	name, ok := dt["name"]
	if ok {
		fmt.Println("有值", name)
	}
	delete(dt, "name") // 删除 name 这一项
	// map 也是引用类型

	// make 函数用于 map、slice、channel 的内存分配，new 则用于各种类型的内存分配
	// make 返回的是初始化后的值，new 返回指针

	// 用 for 实现 while
	sum := 1
	for sum < 1000 {
		sum += sum
	}

	sum1 := 1
	for sum1 < 1000 {
		sum1++
	}

	switch {
	case sum > 500:
		fmt.Println("switch 可以省略表达式，默认匹配为 true 的 case")
	case sum >= 200:
		fmt.Println("只有一个 case 会被执行")
	}

	// panic
	// 是一个内建函数，可以中断原有的控制流程，进入一个令人恐慌的流程中。
	// 当函数 F 调用 panic，函数 F 的执行被中断，但是 F 中的延迟函数会正常执行，然后 F 返回到调用它的地方。
	// 在调用的地方，F 的行为就像调用了 panic。
	// 这一过程继续向上，直到发生 panic 的 goroutine 中所有调用的函数返回，此时程序退出。
	// 恐慌可以直接调用 panic 产生。也可以由运行时错误产生，例如访问越界的数组。

	// recover
	// 是一个内建的函数，可以让进入令人恐慌的流程中的 goroutine 恢复过来
	// recover 仅在延迟函数 (defer 指定的函数) 中有效
	// 在正常的执行过程中，调用 recover 会返回 nil，并且没有其它任何效果
	// 如果当前的 goroutine 陷入恐慌，调用 recover 可以捕获到 panic 的输入值，并且恢复正常的执行
	defer func() {
		if x := recover(); x != nil {
			// 相关处理，日志等
		}
	}()

	// 空接口类型断言
	// element.(type) 语法不能在 switch 外的任何逻辑里面使用

	// select
	// select 默认是阻塞的，只有当监听的 channel 中有发送或接收可以进行时才会运行
	// 当多个 channel 都准备好的时候，select 是随机的选择一个执行的
	// 用于同时监控多个 channel，但是只执行一次，如果要监听多次，需要用 for
	// 可以通过 time.After 来控制超时
	c := make(chan int)
	o := make(chan bool)
	go func() {
		for {
			select {
			case v := <-c:
				println(v)
			case <-time.After(5 * time.Second):
				println("timeout")
				o <- true
				break
			}
		}
	}()
	<-o

	runtime.Goexit()       // 退出当前 goroutine，但已压栈的 defer 会继续执行
	runtime.Gosched()      // 让出当前 goroutine 的执行权限，调度器安排其它等待的任务运行，并在下次某个时候从该位置恢复执行
	runtime.NumCPU()       // 返回 cpu 核数量
	runtime.NumGoroutine() // 返回正在执行和排队的任务总数
	runtime.GOMAXPROCS(1)  // 用来设置可以并行计算的 cpu 核数的最大值，并返回之前的值
}
