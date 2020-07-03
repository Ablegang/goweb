package app

import (
	"github.com/gin-gonic/gin"
	resp "goweb/pkg/response"
	"time"
)

// 注册路由
func registerRoute(r *gin.Engine) {
	// 跨域处理
	r.Use(Cors())

	// 根目录
	r.GET("/", func(c *gin.Context) {
		resp.String(c, 200, time.Now().Format("2006/01/02 - 15:04:05"))
	})
}
