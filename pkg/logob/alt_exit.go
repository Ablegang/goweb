package logob

import (
	"fmt"
	"os"
)

// 已注册的 exit handlers
// 带顺序的 slice ，runHandlers 时，会按 slice 顺序一一执行 handler
var handlers []func()

func runHandler(handler func()) {
	// 捕获所有 handler 可能出现的 panic
	// 有了 recover ，程序不会退出，会继续进入下一次循环（见 runHandlers）
	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "exit handler 执行失败", err)
		}
	}()

	handler()
}

// 执行所有已注册的 handler
func runHandlers() {
	for _, handler := range handlers {
		runHandler(handler)
	}
}

// 退出
func Exit(code int){
	// 退出之前执行所有已注册的 handler
	runHandlers()
	os.Exit(code)
}

// 注册 exit handler，追加到最后
func RegisterExitHandler(handler func()){
	handlers = append(handlers,handler)
}

// 注册 exit handler，添加到最前
func DeferExitHandler(handler func()){
	handlers = append([]func(){handler},handlers...)
}