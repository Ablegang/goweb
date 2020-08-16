package app

import (
	"github.com/gin-gonic/gin"
	"goweb/app/handlers/blog"
	"goweb/app/handlers/show"
	"goweb/pkg/passport"
)

// 注册路由
func registerRoute(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.String(200, "404")
	})

	// 跨域处理
	r.Use(Cors())

	// blog
	blogRouter(r.Group("api/blog"))

	// show
	showRouter(r.Group("api/show"))
}

// 博客路由
func blogRouter(r *gin.RouterGroup) {
	// CI
	r.Any("ci", blog.CI)
}

// show 路由
func showRouter(r *gin.RouterGroup) {
	// 登录
	r.Any("login", show.Login)

	// 登录组
	auth := r.Group("").Use(passport.JwtAuth())
	{
		// 用户信息
		auth.Any("user.info", show.UserInfo)
	}

	// vip 组
	vip := r.Group("").Use(passport.JwtAuth(), Vip())
	{
		// 标的列表（观察中、往期）
		vip.Any("quote.ls", show.QuoteList)
	}

	// 超管组
	superAdmin := r.Group("").Use(passport.JwtAuth(), SuperAdmin())
	{
		// 添加账号
		superAdmin.Any("user.add", show.AddUser)
		// 添加标的
		superAdmin.Any("quote.add", show.QuoteAdd)
		// 删除标的
		superAdmin.Any("quote.del", show.QuoteDel)
		// 标的信息
		superAdmin.Any("quote.info", show.QuoteInfo)
		// 下架标的
		superAdmin.Any("quote.off", show.QuoteOff)
	}
}
