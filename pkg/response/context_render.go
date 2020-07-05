// 此包用来重写 gin.Context 的响应方法，以及所有响应相关的函数都封装在这
// 主要用于做一些除了响应之外的其它操作
// 比如目前就通过 c.Set 来给 logs.RequestAndResponseLog 提供记录 Response 内容的支撑
// 此文件只写关于 gin.Context 响应类的操作

package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goweb/pkg/helper"
	"goweb/pkg/hot"
	"time"
)

var Rely = 1

const (
	DefaultFailCode    = -1
	DefaultSuccessCode = 0
	DefaultFailMsg     = "failed"
	DefaultSuccessMsg  = "ok"
)

// 响应 Json
func JSON(c *gin.Context, code int, g gin.H) {
	// set 是为了让记录日志的中间件能够取到数据
	c.Set("responseCode", code)
	c.Set("responseBody", g)
	c.Set("responseHeader", c.Writer.Header())
	c.JSON(code, g)
}

// 响应 String
func String(c *gin.Context, code int, format string, values ...interface{}) {

	// set 是为了让记录日志的中间件能够取到数据
	c.Set("responseCode", code)
	c.Set("responseBody", fmt.Sprintf(format, values...))
	c.Set("responseHeader", c.Writer.Header())

	c.String(code, format, values...)
}

// 成功响应 Json
func SuccessJson(c *gin.Context, g gin.H) {
	JSON(c, 200, gin.H{
		"code":         DefaultSuccessCode,
		"msg":          DefaultSuccessMsg,
		"api":          c.Request.Host + c.Request.URL.Path + "?" + c.Request.URL.RawQuery,
		"responseTime": time.Now().Format(hot.GetTimeCommonFormat()),
		"data":         g,
	})
}

// 失败响应 Json
func FailJson(c *gin.Context, g gin.H, values ...interface{}) {
	code := DefaultFailCode
	if helper.IssetArrayIndex(values, 0) {
		code, _ = values[0].(int)
	}

	msg := DefaultFailMsg
	if helper.IssetArrayIndex(values, 1) {
		msg, _ = values[1].(string)
	}

	JSON(c, 400, gin.H{
		"code":         code,
		"msg":          msg,
		"api":          c.Request.Host + c.Request.URL.Path + "?" + c.Request.URL.RawQuery,
		"responseTime": time.Now().Format(hot.GetTimeCommonFormat()),
		"data":         g,
	})
}
