// 业务型中间件写在此文件下

package app

import (
	"github.com/gin-gonic/gin"
	"goweb/pkg/hot"
	resp "goweb/pkg/response"
)

// 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin, _ := hot.GetConfig("ResponseHeader.CorsOrigin").(string)
		credentials, _ := hot.GetConfig("ResponseHeader.CorsCredentials").(string)
		author, _ := hot.GetConfig("ResponseHeader.Author").(string)
		powerBy, _ := hot.GetConfig("ResponseHeader.PowerBy").(string)
		powerBy2, _ := hot.GetConfig("ResponseHeader.PowerBy2").(string)

		// 响应头要前置处理，才能被下一个 context 继承
		// 否则会先处理 Next ，则 Next 的 context 并未写入响应头
		// 详情可看 c.Next 方法
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", credentials)
		c.Header("Author", author)
		c.Header("PowerBy", powerBy)
		c.Header("PowerBy2", powerBy2)
		c.Next()
	}
}

// jwt 验证中间件
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if len(token) <= 25 {
			resp.FailJson(c, "请先登录！")
			c.Abort()
			return
		}



		c.Next()
	}
}
