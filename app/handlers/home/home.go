package home

import (
	"github.com/gin-gonic/gin"
	"goweb/app/helper"
)

func Index(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code":     0,
		"msg":      "ok",
		"minLevel": helper.Get("user.ok"),
	})
}
