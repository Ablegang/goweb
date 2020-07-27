package show

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"goweb/app/models"
	"goweb/app/models/show"
	"goweb/pkg/dingrobot"
	"goweb/pkg/quotes"
	"goweb/pkg/request"
	resp "goweb/pkg/response"
	"os"
	"time"
)

// 基础表单字段 reason
type QuoteReason struct {
	Reason string `json:"reason" form:"reason" validate:"min=10,max=255"`
}

// 基础表单字段 id
type QuoteId struct {
	Id int `json:"id" form:"id" validate:"numeric"`
}

// 添加表单
type AddQuoteForm struct {
	Key string `json:"key" form:"key" validate:"len=7"`
	QuoteReason
}

// 编辑表单
type EditQuoteForm struct {
	QuoteId
	AddQuoteForm
}

// 下架表单
type OffQuoteForm struct {
	QuoteId
	QuoteReason
}

// 标的状态常量
const (
	DisplayState = "display"
	OffState     = "off"
)

// 钉钉通知相关...
var (
	RobotToken  = os.Getenv("QUOTES_DING_ACCESS_TOKEN")
	AtMobile    = os.Getenv("AT_MOBILE")
	AddTemplate = "# 新标入选 \n %s（%s） \n > 入选理由：%s \n @%s"
	AddTitle    = "新标入选通知"
)

// 添加标的
func QuoteAdd(c *gin.Context) {
	// 入参
	req := &AddQuoteForm{}
	if err := request.Bind(c, req); err != nil {
		resp.FailJson(c, gin.H{}, -1, err.Error())
		return
	}

	// 唯一性校验
	has, _ := models.Show().Exist(&show.Quote{
		Key:    req.Key,
		Status: DisplayState,
	})
	if has {
		resp.FailJson(c, gin.H{}, -1, "该标的已在展示中")
		return
	}

	// 取标的信息
	driver := quotes.New(quotes.WyResource)
	driver.SetKeys([]string{req.Key})
	data, err := driver.GetQuotes()
	if err != nil {
		resp.FailJson(c, gin.H{}, -1, "抓取数据失败")
		logrus.Errorln("抓取数据失败", err)
		return
	}
	if len(data) == 0 {
		resp.FailJson(c, gin.H{}, -1, "不存在的标的标识")
		return
	}

	// 存储信息
	_, err = models.Show().Insert(&show.Quote{
		Name:         data[0].Name,
		Number:       data[0].Number,
		Key:          req.Key,
		InitialPrice: data[0].NowPrice,
		TodayPrice:   data[0].NowPrice,
		AddReason:    req.Reason,
		OffReason:    "",
		Status:       DisplayState,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		resp.FailJson(c, gin.H{}, -1, err.Error())
		return
	}

	// 钉钉通知
	dingrobot.Markdown(&dingrobot.MarkdownParams{
		Ac:      RobotToken,
		Md:      fmt.Sprintf(AddTemplate, data[0].Name, data[0].Number, req.Reason, AtMobile),
		Title:   AddTitle,
		At:      []string{AtMobile},
		IsAtAll: false,
	})

	resp.SuccessJson(c, gin.H{})
	return
}

// 编辑标的
func QuoteEdit(c *gin.Context) {

}

// 删除标的
func QuoteDel(c *gin.Context) {

}

// 标的信息
func QuoteInfo(c *gin.Context) {

}

// 标的列表
func QuoteList(c *gin.Context) {

}

// 下架标的
func QuoteOff(c *gin.Context) {

}
