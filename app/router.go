package app

import (
	"github.com/gin-gonic/gin"
	"goweb/app/handlers/blog"
	"goweb/app/handlers/show"
	"goweb/pkg/passport"
)

// 注册路由
func registerRoute(r *gin.Engine) {
	// 跨域处理
	r.Use(Cors())

	// blog
	blogRouter(r.Group("blog"))

	// show
	showRouter(r.Group("show"))
}

// 博客路由
func blogRouter(r *gin.RouterGroup) {
	// 分类列表
	r.GET("ls", blog.Ls)
	// 文章列表
	r.GET("posts", blog.Posts)
	// 文章内容
	r.GET("posts.detail", blog.PostsDetail)
	// CI
	r.POST("ci", blog.CI)
}

// show 路由
func showRouter(r *gin.RouterGroup) {
	// 登录
	r.Any("login", show.Login)

	auth := r.Use(passport.JwtAuth())

	// 添加账号
	auth.Any("admin.add", show.AddAdmin)
	// 添加标的
	auth.Any("quote.add", show.QuoteAdd)
	// 删除标的
	auth.Any("quote.del", show.QuoteDel)
	// 标的信息
	auth.Any("quote.info", show.QuoteInfo)
	// 标的列表（观察中、往期）
	auth.Any("quote.ls", show.QuoteList)
	// 下架标的
	auth.Any("quote.off", show.QuoteOff)
}
