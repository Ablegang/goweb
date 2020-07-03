// 此包用来重写 gin.Context 的响应方法
// 主要用于做一些除了响应之外的其它操作
// 比如目前就通过 c.Set 来给 logs.RequestAndResponseLog 提供记录 Response 内容的支撑
// 此文件只写关于 gin.Context 响应类的操作

package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var Rely = 1

// 响应 Json
func Json(c *gin.Context, code int, g gin.H) {
	// set 是为了让记录日志的中间件能够取到数据
	c.Set("responseCode", code)
	c.Set("responseBody", g)
	c.Set("responseHeader", c.Writer.Header())
	c.JSON(code, g)
	return
}

// 响应 String
func String(c *gin.Context, code int, format string, values ...interface{}) {
	
	// set 是为了让记录日志的中间件能够取到数据
	c.Set("responseCode", code)
	c.Set("responseBody", fmt.Sprintf(format, values...))
	c.Set("responseHeader", c.Writer.Header())

	c.String(code, format, values...)
	return
}

// 更多方法，需要时，再根据 gin.Context 定义
