// 框架入口

package app

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"goweb/pkg/logs"
	"goweb/pkg/logs/loghooks"
	"goweb/pkg/response"
	"io"
	"os"
)

// gin 路由器实例
var r *gin.Engine

// 初始化框架
func init() {
	// 环境变量
	loadEnv()
	// 注册日志组件
	registerLogger()
	// 实例化路由器
	r = router()
	// 注册路由
	registerRoute(r)
}

// 框架启动
func Start() {
	port := os.Getenv("PORT")
	err := r.Run(port)
	if err != nil {
		log.Panicln("启动失败：", err)
	}
}

// 注册日志组件
func registerLogger() {
	// 设置最小 log 级别
	log.SetLevel(log.TraceLevel)
	// 设置取堆栈信息
	log.SetReportCaller(true)

	// 注册 Hooks...
	log.AddHook(loghooks.NewEmailNotify())
	log.AddHook(loghooks.NewFileWriter())
}

// 载入 .env
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln("env load failed：", err)
	}
}

// 实例化 gin 路由器
func router() *gin.Engine {
	// 实例化路由器
	r = gin.New()

	// recovery 相关处理：记录 panic 日志到 file
	recoverWriter := io.MultiWriter(gin.DefaultErrorWriter, &logs.CustomFileWriter{
		LogMode:          "daily",
		Dir:              "storage/logs/ginErr/",
		FileNameFormater: "2006-01-02.txt",
		Perm:             os.FileMode(0777),
	})
	r.Use(response.RecoveryWithWriter(recoverWriter))

	// logger 相关处理：记录默认的 gin 请求日志到 file
	logWriter := io.MultiWriter(gin.DefaultWriter, logs.NewCustomFileWriter())
	r.Use(gin.LoggerWithWriter(logWriter))

	// requests 和 responses 记录到 file
	r.Use(logs.RequestAndResponseLog())

	return r
}
