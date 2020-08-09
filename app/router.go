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
	r.Any("ls", blog.Ls)
	// 文章列表
	r.Any("posts", blog.Posts)
	// 文章内容
	r.Any("posts.detail", blog.PostsDetail)
	// CI
	r.POST("ci", blog.CI)
}

// show 路由
func showRouter(r *gin.RouterGroup) {
	// 登录
	r.Any("login", show.Login)

	// 登录用户
	auth := r.Use(passport.JwtAuth())
	// vip 用户
	vip := auth.Use(Vip())
	// 超管用户
	superAdmin := auth.Use(SuperAdmin())

	// 用户信息
	auth.Any("user.info", show.UserInfo)

	// 标的列表（观察中、往期）
	vip.Any("quote.ls", show.QuoteList)

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
