// 业务型中间件写在此文件下

package app

import (
	"github.com/gin-gonic/gin"
	"goweb/app/models"
	"goweb/app/models/show"
	"goweb/pkg/hot"
	"goweb/pkg/response"
	"net/http"
)

// 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin, _ := hot.GetConfig("ResponseHeader.CorsOrigin").(string)
		credentials, _ := hot.GetConfig("ResponseHeader.CorsCredentials").(string)
		author, _ := hot.GetConfig("ResponseHeader.Author").(string)
		powerBy, _ := hot.GetConfig("ResponseHeader.PowerBy").(string)
		powerBy2, _ := hot.GetConfig("ResponseHeader.PowerBy2").(string)
		headers, _ := hot.GetConfig("ResponseHeader.Headers").(string)

		// 响应头要前置处理，才能被下一个 context 继承
		// 否则会先处理 Next ，则 Next 的 context 并未写入响应头
		// 详情可看 c.Next 方法
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", credentials)
		c.Header("Access-Control-Allow-Headers", headers)
		c.Header("Author", author)
		c.Header("PowerBy", powerBy)
		c.Header("PowerBy2", powerBy2)

		// 如果是 OPTIONS 请求，直接返回 200 状态码，以便于前端跨域
		if c.Request.Method == "OPTIONS" {
			c.String(http.StatusOK, "")
			c.Abort()
		}

		c.Next()
	}
}

// 会员过滤
func Vip() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 超管拥有所有权限
		CheckRole(c, []int64{show.ADMIN, show.VIP}, -2, "您不是 vip")
	}
}

// 超管过滤
func SuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		CheckRole(c, []int64{show.ADMIN}, -3, "您不是超管")
	}
}

// 检查用户组
func CheckRole(c *gin.Context, roleIds []int64, code int, msg string) {
	userI, _ := c.Get("AuthUser")
	if userI == nil {
		response.FailJson(c, gin.H{}, code, msg)
		c.Abort()
		return
	}
	user, _ := userI.(map[string]interface{})

	query := models.Show().Where("user_id = ?", user["Id"])
	total, err := query.In("role_id", roleIds).Count(&show.UserRole{})

	if err != nil || total <= 0 {
		response.FailJson(c, gin.H{}, code, msg)
		c.Abort()
		return
	}

	c.Next()
}
