package dingrobot

import (
	"github.com/sirupsen/logrus"
)

// 钉钉通知 err 不可使用 logrus.Error 及以上级别，否则会陷入死循环：失败--触发钩子--钩子调推送--失败
// 所以只有用不会触发钩子的错误级别

// Markdown 参数
type MarkdownParams struct {
	Ac         string
	Md         string
	Title      string
	At         []string
	IsAtAll    bool
	ErrHandler func(err error)
}

// Markdown 通知
func Markdown(params *MarkdownParams) {
	robot := NewRobot(params.Ac)
	msg := NewMessageBuilder(TypeMarkdown).Markdown(params.Title, params.Md).At(params.At, params.IsAtAll).Build()
	err := robot.SendMessage(msg)
	if err != nil {
		if params.ErrHandler != nil {
			params.ErrHandler(err)
		} else {
			logrus.Println("钉钉通知失败", err)
		}
	}
}

// 文本通知参数
type TextParams struct {
	Ac         string
	Text       string
	At         []string
	IsAtAll    bool
	ErrHandler func(err error)
}

// 文本通知
func Text(params *TextParams) {
	robot := NewRobot(params.Ac)
	msg := NewMessageBuilder(TypeText).Text(params.Text).At(params.At, params.IsAtAll).Build()
	err := robot.SendMessage(msg)
	if err != nil {
		if params.ErrHandler != nil {
			params.ErrHandler(err)
		} else {
			logrus.Println("钉钉通知失败", err)
		}
	}
}

// Link 通知
func Link(ac string, link map[string]string, at []string, isAtAll bool) {
	robot := NewRobot(ac)
	build := NewMessageBuilder(TypeLink).At(at, isAtAll)
	build = build.Link(link["title"], link["text"], link["msgUrl"], link["picUrl"])
	msg := build.Build()
	err := robot.SendMessage(msg)
	if err != nil {
		logrus.Println("钉钉通知失败", err)
	}
}

// ActionCard 通知
// https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq
func ActionCard() {

}

// FeedCard 通知
// 详情看文档
func FeedCard() {

}
