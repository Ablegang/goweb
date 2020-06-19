// 框架入口

package app

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"goweb/pkg/logs"
	"io"
)

var r *gin.Engine

func init() {
	registerLogger()

	loadEnv()

	r = gin.New()

	r.Use(gin.Recovery())

	// 自定义 gin logger ，主要为了将日志数据按日输出到自定义的目录内
	writer := io.MultiWriter(gin.DefaultWriter, logs.NewRequestWriter())
	r.Use(gin.LoggerWithWriter(writer))

	// 注册路由
	registerRoute(r)

	return
}

// 框架启动
func Start() {
	port, _ := Get("app", "port").(string)
	err := r.Run(port)
	if err != nil {
		log.Panicln("启动失败：", err)
	}
}

func registerLogger() {
	// 直接使用 log 包的函数，包内部会实例化一个 std 作为通用 logger
	// std 是指针类型，且是包变量，程序运行时就已经注册

	// 设置最小 log 级别
	level, _ := Get("log", "minLevel").(log.Level)
	log.SetLevel(level)

	// 设置取堆栈信息
	reportCaller, _ := Get("log", "reportCaller").(bool)
	log.SetReportCaller(reportCaller)

	// 注册 Hooks...
	hooks, _ := Get("log", "hooks").([]log.Hook)
	for _, hook := range hooks {
		log.AddHook(hook)
	}

	return
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln("env load failed：", err)
	}
	return
}
