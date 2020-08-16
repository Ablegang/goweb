package blog

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/webhooks.v5/github"
	"goweb/pkg/dingrobot"
	"goweb/pkg/helper"
	"os"
	"os/exec"
	"time"
)

// CI
func CI(c *gin.Context) {

	// hook 校验
	if os.Getenv("GITEE_WEBHOOK_BLOG_SECRET") != c.Request.Header.Get("X-Gitee-Token") {
		logrus.Errorln("BLOG CI 失败 ：密码错误")
		c.String(404, "404 not found")
		return
	}

	// 请求体
	var payload GiteePayload
	if err := c.ShouldBind(&payload); err != nil {
		logrus.Errorln("BLOG CI 失败 ：", err)
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
	md += "- 更新内容：" + payload.HeadCommit.Message + "\n"
	md += "- 作者：" + payload.HeadCommit.Committer.Name + "\n"
	md += "- Email：" + payload.HeadCommit.Committer.Email + "\n"
	md += "- 社区名：" + payload.HeadCommit.Committer.Username + "\n"
	dingrobot.Markdown(&dingrobot.MarkdownParams{
		Ac:      os.Getenv("BLOG_DING_ACCESS_TOKEN"),
		Md:      md,
		Title:   "Blog CI/CD",
		At:      []string{},
		IsAtAll: false,
	})
	return
}

// github ci
func CI2(c *gin.Context) {
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

	c.JSON(200, data)
}

// gitee payload
type GiteePayload struct {
	Ref                string     `json:"ref"`
	Before             string     `json:"before"`
	After              string     `json:"after"`
	Created            bool       `json:"created"`
	Deleted            bool       `json:"deleted"`
	Compare            string     `json:"compare"`
	Commits            []Commits  `json:"commits"`
	HeadCommit         HeadCommit `json:"head_commit"`
	TotalCommitsCount  int        `json:"total_commits_count"`
	CommitsMoreThanTen bool       `json:"commits_more_than_ten"`
	Repository         Repository `json:"repository"`
	Project            Project    `json:"project"`
	UserID             int        `json:"user_id"`
	UserName           string     `json:"user_name"`
	User               User       `json:"user"`
	Pusher             Pusher     `json:"pusher"`
	Sender             Sender     `json:"sender"`
	Enterprise         Enterprise `json:"enterprise"`
	HookName           string     `json:"hook_name"`
	HookID             int        `json:"hook_id"`
	HookURL            string     `json:"hook_url"`
	Password           string     `json:"password"`
	Timestamp          string     `json:"timestamp"`
	Sign               string     `json:"sign"`
}
type Author struct {
	Time     time.Time `json:"time"`
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	UserName string    `json:"user_name"`
	URL      string    `json:"url"`
}
type Committer struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	UserName string `json:"user_name"`
	URL      string `json:"url"`
}
type Commits struct {
	ID        string      `json:"id"`
	TreeID    string      `json:"tree_id"`
	ParentIds []string    `json:"parent_ids"`
	Distinct  bool        `json:"distinct"`
	Message   string      `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
	URL       string      `json:"url"`
	Author    Author      `json:"author"`
	Committer Committer   `json:"committer"`
	Added     interface{} `json:"added"`
	Removed   interface{} `json:"removed"`
	Modified  []string    `json:"modified"`
}
type HeadCommit struct {
	ID        string      `json:"id"`
	TreeID    string      `json:"tree_id"`
	ParentIds []string    `json:"parent_ids"`
	Distinct  bool        `json:"distinct"`
	Message   string      `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
	URL       string      `json:"url"`
	Author    Author      `json:"author"`
	Committer Committer   `json:"committer"`
	Added     interface{} `json:"added"`
	Removed   interface{} `json:"removed"`
	Modified  []string    `json:"modified"`
}
type Owner struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	UserName  string `json:"user_name"`
	URL       string `json:"url"`
}
type Repository struct {
	ID                int         `json:"id"`
	Name              string      `json:"name"`
	Path              string      `json:"path"`
	FullName          string      `json:"full_name"`
	Owner             Owner       `json:"owner"`
	Private           bool        `json:"private"`
	HTMLURL           string      `json:"html_url"`
	URL               string      `json:"url"`
	Description       string      `json:"description"`
	Fork              bool        `json:"fork"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	PushedAt          time.Time   `json:"pushed_at"`
	GitURL            string      `json:"git_url"`
	SSHURL            string      `json:"ssh_url"`
	CloneURL          string      `json:"clone_url"`
	SvnURL            string      `json:"svn_url"`
	GitHTTPURL        string      `json:"git_http_url"`
	GitSSHURL         string      `json:"git_ssh_url"`
	GitSvnURL         string      `json:"git_svn_url"`
	Homepage          interface{} `json:"homepage"`
	StargazersCount   int         `json:"stargazers_count"`
	WatchersCount     int         `json:"watchers_count"`
	ForksCount        int         `json:"forks_count"`
	Language          string      `json:"language"`
	HasIssues         bool        `json:"has_issues"`
	HasWiki           bool        `json:"has_wiki"`
	HasPages          bool        `json:"has_pages"`
	License           interface{} `json:"license"`
	OpenIssuesCount   int         `json:"open_issues_count"`
	DefaultBranch     string      `json:"default_branch"`
	Namespace         string      `json:"namespace"`
	NameWithNamespace string      `json:"name_with_namespace"`
	PathWithNamespace string      `json:"path_with_namespace"`
}
type Project struct {
	ID                int         `json:"id"`
	Name              string      `json:"name"`
	Path              string      `json:"path"`
	FullName          string      `json:"full_name"`
	Owner             Owner       `json:"owner"`
	Private           bool        `json:"private"`
	HTMLURL           string      `json:"html_url"`
	URL               string      `json:"url"`
	Description       string      `json:"description"`
	Fork              bool        `json:"fork"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	PushedAt          time.Time   `json:"pushed_at"`
	GitURL            string      `json:"git_url"`
	SSHURL            string      `json:"ssh_url"`
	CloneURL          string      `json:"clone_url"`
	SvnURL            string      `json:"svn_url"`
	GitHTTPURL        string      `json:"git_http_url"`
	GitSSHURL         string      `json:"git_ssh_url"`
	GitSvnURL         string      `json:"git_svn_url"`
	Homepage          interface{} `json:"homepage"`
	StargazersCount   int         `json:"stargazers_count"`
	WatchersCount     int         `json:"watchers_count"`
	ForksCount        int         `json:"forks_count"`
	Language          string      `json:"language"`
	HasIssues         bool        `json:"has_issues"`
	HasWiki           bool        `json:"has_wiki"`
	HasPages          bool        `json:"has_pages"`
	License           interface{} `json:"license"`
	OpenIssuesCount   int         `json:"open_issues_count"`
	DefaultBranch     string      `json:"default_branch"`
	Namespace         string      `json:"namespace"`
	NameWithNamespace string      `json:"name_with_namespace"`
	PathWithNamespace string      `json:"path_with_namespace"`
}
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	UserName string `json:"user_name"`
	URL      string `json:"url"`
}
type Pusher struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	UserName string `json:"user_name"`
	URL      string `json:"url"`
}
type Sender struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	UserName  string `json:"user_name"`
	URL       string `json:"url"`
}
type Enterprise struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
