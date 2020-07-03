// 业务型中间件写在此文件下

package app

import (
	"github.com/gin-gonic/gin"
	"goweb/app/helper"
)

// 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin, _ := helper.Get("ResponseHeader.CorsOrigin").(string)
		credentials, _ := helper.Get("ResponseHeader.CorsCredentials").(string)
		author, _ := helper.Get("ResponseHeader.Author").(string)
		powerBy, _ := helper.Get("ResponseHeader.PowerBy").(string)
		powerBy2, _ := helper.Get("ResponseHeader.PowerBy2").(string)

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
