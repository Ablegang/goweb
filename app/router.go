package app

import (
	"github.com/gin-gonic/gin"
	"goweb/pkg/response"
	"time"
)

// 注册路由
func registerRoute(r *gin.Engine) {
	// 跨域处理
	r.Use(Cors())

	r.GET("/", func(c *gin.Context) {
		response.String(c, 200, time.Now().Format("2006/01/01 - 15:04:05"))
	})
}
