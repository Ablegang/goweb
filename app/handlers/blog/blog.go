package blog

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/webhooks.v5/github"
	"goweb/pkg/dingrobot"
	"goweb/pkg/helper"
	"goweb/pkg/hot"
	resp "goweb/pkg/response"
	"io/ioutil"
	"os"
	"os/exec"
)

// 分类列表
func Ls(c *gin.Context) {
	path, rootDir, err := getPath(c)
	if err != nil {
		resp.FailJson(c, gin.H{})
		return
	}

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
	path, rootDir, err := getPath(c)
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
	for _, v := range posts {
		if !v.IsDir() {
			p = append(p, map[string]interface{}{
				"name": v.Name(),
				"ut":   v.ModTime().Format(hot.GetTimeCommonFormat()),
				"path": helper.FormatPath(path, rootDir),
			})
		}
	}

	resp.SuccessJson(c, gin.H{"list": p})
	return
}

// 文章详情
func PostsDetail(c *gin.Context) {
	// 路径校验
	path, rootDir, err := getPath(c)
	if err != nil {
		resp.FailJson(c, gin.H{}, resp.DefaultFailCode, "未指定路径")
		return
	}

	// 入参校验
	name, ok := c.GetQuery("name")
	if !ok {
		resp.FailJson(c, gin.H{}, resp.DefaultFailCode, "请指定文件")
		return
	}

	// 读取文件
	content, err := ioutil.ReadFile(path + name)
	if err != nil {
		notice := "网络超时"
		if os.IsNotExist(err) {
			notice = "数据不存在"
		} else {
			logrus.Errorln("博客系统文件系统告警", err, "filePath:"+path+name)
		}

		resp.FailJson(c, gin.H{}, resp.DefaultFailCode, notice)
		return
	}

	file, _ := os.Stat(path + name)

	resp.SuccessJson(c, gin.H{
		"name":    name,
		"path":    helper.FormatPath(path, rootDir),
		"content": string(content),
		"ut":      file.ModTime().Format(hot.GetTimeCommonFormat()),
	})
	return
}

// CI
func CI(c *gin.Context) {
	// hook 校验
	hook, _ := github.New(github.Options.Secret(os.Getenv("GITHUB_WEBHOOK_BLOG_SECRET")))
	payload, err := hook.Parse(c.Request, github.PushEvent, github.PingEvent)
	if err != nil {
		logrus.Errorln("BLOG CI 失败", err)
		c.String(404, "404 not found")
		return
	}

	// payload 解析
	data, ok := payload.(github.PushPayload)
	if !ok {
		logrus.Errorln("BLOG CI 失败", err)
		c.String(404, "404 not found")
		return
	}

	// 执行部署
	rootDir, _ := helper.GetBlogRoot()
	shell := "cd " + "./" + rootDir
	shell += "&& git pull origin master"
	cmd := exec.Command("/bin/sh", "-c", shell)
	_, _ = cmd.Output()

	// 通知
	robot := dingrobot.NewRobot(os.Getenv("LOG_DING_ACCESS_TOKEN"))
	md := "# BLOG 发布了新的内容 \n"
	md += "- 更新内容：" + data.HeadCommit.Message + "\n"
	md += "- 作者：" + data.HeadCommit.Committer.Name + "\n"
	md += "- Email：" + data.HeadCommit.Committer.Email + "\n"
	md += "- 社区名：" + data.HeadCommit.Committer.Username + "\n"
	msg := dingrobot.NewMessageBuilder(dingrobot.TypeMarkdown).Markdown("市场监控", md).Build()
	err = robot.SendMessage(msg)
	if err != nil {
		logrus.Errorln("钉钉行情推送失败", err)
	}
	return
}

// 获取 Path
func getPath(c *gin.Context) (string, string, error) {
	// 根目录配置
	rootDir, _ := helper.GetBlogRoot()

	// 请求参数
	fullPath := rootDir
	path, ok := c.GetQuery("path")
	if ok && path != ".git" {
		fullPath += path
		if string(fullPath[len(fullPath)-1:]) != "/" {
			fullPath += "/"
		}
	}

	return fullPath, rootDir, nil
}
