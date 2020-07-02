package app

import (
	"github.com/gin-gonic/gin"
)

func registerRoute(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.String(200,"It works!")
	})
}