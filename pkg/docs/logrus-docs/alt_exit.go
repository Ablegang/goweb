package logrus_docs

import (
	"fmt"
	"os"
)

// 这里注册的 handlers 是包空间的，不论是哪个实例，哪个 goroutine ，都将生效
var handlers []func()

func runHandler(handler func()) {
	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error: Logrus exit handler error:", err)
		}
	}()

	handler()
}

// 执行所有 ExitHandler
func runHandlers() {
	for _, handler := range handlers {
		runHandler(handler)
	}
}

// 退出函数
func Exit(code int) {
	runHandlers()
	os.Exit(code)
}

// 从栈尾注册一个 ExitHandler
func RegisterExitHandler(handler func()) {
	handlers = append(handlers, handler)
}

// 从栈首注册一个 ExitHandler
func DeferExitHandler(handler func()) {
	handlers = append([]func(){handler}, handlers...)
}
