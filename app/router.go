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
	r.Any("lg", show.Login)

	// 需登录
	auth := r.Use(passport.JwtAuth())
	auth.Any("add", show.Add)
}
