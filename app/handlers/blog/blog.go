package blog

import (
	"github.com/gin-gonic/gin"
	"goweb/pkg/helper"
	"goweb/pkg/hot"
	resp "goweb/pkg/response"
	"io/ioutil"
)

// 分类列表
func Ls(c *gin.Context) {
	path, err := getPath(c)
	if err != nil {
		resp.FailJson(c, gin.H{})
		return
	}

	rootDir, _ := helper.GetBlogRoot()

	// 取数据
	dirs, err := helper.RecursiveGetDirList(path, rootDir)
	if err != nil {
		// path 参数由外部传递，所以 err 时，不需要告警
		resp.FailJson(c, gin.H{}, resp.DefaultFailCode, "路径非法")
		return
	}

	resp.SuccessJson(c, gin.H{"list": dirs})
	return
}

// 文章列表
func Posts(c *gin.Context) {
	path, err := getPath(c)
	if err != nil {
		resp.FailJson(c, gin.H{}, resp.DefaultFailCode, "未指定路径")
		return
	}

	posts, err := ioutil.ReadDir(path)
	if err != nil {
		resp.FailJson(c, gin.H{}, resp.DefaultFailCode, "路径非法")
		return
	}

	p := make([]map[string]interface{}, 0)
	except, _ := hot.GetConfig("blog.root").(string)
	for _, v := range posts {
		if !v.IsDir() {
			p = append(p, map[string]interface{}{
				"name": v.Name(),
				"ut":   v.ModTime().Format(hot.GetTimeCommonFormat()),
				"path": helper.FormatPath(path, except),
			})
		}
	}

	resp.SuccessJson(c, gin.H{"list": p})
	return
}

// 文章详情
func PostsDetail(c *gin.Context) {

}

// 获取 Path
func getPath(c *gin.Context) (string, error) {
	// 根目录配置
	rootDir, _ := helper.GetBlogRoot()

	// 请求参数
	path, ok := c.GetQuery("path")
	if ok {
		rootDir += path
		if string(rootDir[len(rootDir)-1:]) != "/" {
			rootDir += "/"
		}
	}

	return rootDir, nil
}
