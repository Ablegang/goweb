package app

import (
	"github.com/gin-gonic/gin"
	"goweb/app/handlers/home"
)

func registerRoute(r *gin.Engine) {
	r.Any("/",home.Index)

	r.
}