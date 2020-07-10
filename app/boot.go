// 框架入口

package app

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"goweb/app/models"

	// 环境变量提前载入，要比框架更早
	_ "goweb/pkg/env"
	"goweb/pkg/logs"
	"goweb/pkg/logs/loghooks"
	resp "goweb/pkg/response"
	"io"
	"os"
)

// gin 路由器实例
var r *gin.Engine

// 初始化框架
func init() {
	// 注册日志组件
	registerLogger()
	// 实例化路由器
	r = router()
	// 注册路由
	registerRoute(r)
	// 注册数据库引擎
	models.Load()
}

// 框架启动
func Start() {
	// 启动 cron 进程
	go startCron()

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
	log.AddHook(loghooks.NewDingNotify(os.Getenv("LOG_DING_ACCESS_TOKEN")))
}

// 实例化 gin 路由器
func router() *gin.Engine {
	// 实例化路由器
	r = gin.New()

	// 全局唯一 ID
	r.Use(resp.RequestId())

	// recovery 相关处理：记录 panic 日志到 file
	recoverWriter := io.MultiWriter(gin.DefaultErrorWriter, &logs.CustomFileWriter{
		LogMode:          os.Getenv("RECOVER_LOG_MODE"),
		Dir:              "storage/" + os.Getenv("RECOVER_LOG_DIR"),
		FileNameFormater: os.Getenv("RECOVER_LOG_FILEFORMATER"),
		Perm:             os.FileMode(0777),
		IsDingRobot:      true,
	})
	r.Use(resp.RecoveryWithWriter(recoverWriter))

	// logger 相关处理：记录默认的 gin 请求日志到 file
	logWriter := io.MultiWriter(gin.DefaultWriter, &logs.CustomFileWriter{
		LogMode:          os.Getenv("GIN_STD_LOG_MODE"),
		Dir:              "storage/" + os.Getenv("GIN_STD_LOG_DIR"),
		FileNameFormater: os.Getenv("GIN_STD_LOG_FILEFORMATER"),
		Perm:             os.FileMode(0777),
		IsDingRobot:      false,
	})
	r.Use(gin.LoggerWithWriter(logWriter))

	// requests 和 responses 记录到 file
	r.Use(logs.RequestAndResponseLog(&logs.CustomFileWriter{
		LogMode:          os.Getenv("REQUEST_LOG_MODE"),
		Dir:              "storage/" + os.Getenv("REQUEST_LOG_DIR"),
		FileNameFormater: os.Getenv("REQUEST_LOG_FILEFORMATER"),
		Perm:             os.FileMode(0777),
		IsDingRobot:      false,
	}))

	return r
}
