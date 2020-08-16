package blog

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/webhooks.v5/github"
	"goweb/pkg/dingrobot"
	"goweb/pkg/helper"
	"os"
	"os/exec"
)

// CI
func CI(c *gin.Context) {
	c.JSON(200, gin.H{
		"token":   c.Request.Header.Get("X-Gitee-Token"),
		"payload": c.Request.PostForm,
	})
	return

	// hook 校验
	hook, _ := github.New(github.Options.Secret(os.Getenv("GITEE_WEBHOOK_BLOG_SECRET")))
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
	md := "# BLOG 发布了新的内容 \n"
	md += "- 更新内容：" + data.HeadCommit.Message + "\n"
	md += "- 作者：" + data.HeadCommit.Committer.Name + "\n"
	md += "- Email：" + data.HeadCommit.Committer.Email + "\n"
	md += "- 社区名：" + data.HeadCommit.Committer.Username + "\n"
	dingrobot.Markdown(&dingrobot.MarkdownParams{
		Ac:      os.Getenv("BLOG_DING_ACCESS_TOKEN"),
		Md:      md,
		Title:   "Blog CI/CD",
		At:      []string{},
		IsAtAll: true,
	})
	return
}
