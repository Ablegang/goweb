package app

import "github.com/gin-gonic/gin"

func registerRoute(r *gin.Engine) {
	r.Any("/")
}