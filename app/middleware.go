// 业务型中间件写在此文件下

package app

import (
	"github.com/gin-gonic/gin"
)

// 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 响应头要前置处理，才能被下一个 context 继承
		// 否则会先处理 Next ，则 Next 的 context 并未写入响应头
		// 详情可看 c.Next 方法
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Author", "Object(Ablegang)")
		c.Header("PowerBy", "www.goenv.cn / www.goenv.net")
		c.Next()
	}
}
