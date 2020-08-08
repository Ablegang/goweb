package show

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"goweb/app/models"
	"goweb/app/models/show"
	"goweb/pkg/dingrobot"
	"goweb/pkg/helper"
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
		CreatedAt:    helper.JsonTime(time.Now()),
		UpdatedAt:    helper.JsonTime(time.Now()),
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

// 删除标的
func QuoteDel(c *gin.Context) {
	req := &QuoteId{}
	if err := request.Bind(c, req); err != nil {
		resp.FailJson(c, gin.H{}, -1, err.Error())
		return
	}

	// 查询数据
	quote := show.Quote{}
	has, _ := models.Show().Where("id = ?", req.Id).Get(&quote)
	if !has {
		resp.FailJson(c, gin.H{}, -1, "数据不存在")
		return
	}

	// 删除数据
	res, _ := models.Show().Delete(&quote)
	if res <= 0 {
		resp.FailJson(c, gin.H{}, -1, "删除失败")
		return
	}

	resp.SuccessJson(c, gin.H{})
}

// 标的信息
func QuoteInfo(c *gin.Context) {
	// 入参
	req := &QuoteId{}
	if err := request.Bind(c, req); err != nil {
		resp.FailJson(c, gin.H{}, -1, err.Error())
		return
	}

	// 查询数据
	quote := show.Quote{}
	has, _ := models.Show().Where("id = ?", req.Id).Get(&quote)
	if !has {
		resp.FailJson(c, gin.H{}, -1, "数据不存在")
		return
	}

	resp.SuccessJson(c, quote)
	return
}

// 标的列表
func QuoteList(c *gin.Context) {
	// 入参
	req := &helper.CommonPageForm{}
	if err := request.Bind(c, req); err != nil {
		resp.FailJson(c, gin.H{}, -1, err.Error())
		return
	}

	list := make([]show.Quote, 0)

	// 状态筛选，默认为展示中
	status := DisplayState
	if len(c.Param("status")) > 0 {
		status = c.Param("status")
	}
	query := models.Show().Where("status = ?", status)

	// 关键词筛选
	if len(req.Keyword) > 0 {
		query = query.Where("name like '%?%'", req.Keyword).Or("number like '%?%'", req.Keyword)
	}

	// 总数
	total, _ := query.Count(&show.Quote{})

	err := query.OrderBy("created_at desc").Limit(req.Limit, (req.Page-1)*req.Limit).Find(&list)
	if err != nil {
		resp.FailJson(c, gin.H{}, -1, err.Error())
		return
	}

	resp.SuccessJson(c, gin.H{
		"total": total,
		"list":  list,
	})
}

// 下架标的
func QuoteOff(c *gin.Context) {
	req := &OffQuoteForm{}
	if err := request.Bind(c, req); err != nil {
		resp.FailJson(c, gin.H{}, -1, err.Error())
		return
	}

	// 查询数据
	quote := show.Quote{}
	has, _ := models.Show().Where("id = ?", req.Id).Get(&quote)
	if !has {
		resp.FailJson(c, gin.H{}, -1, "数据不存在")
		return
	}

	// 下架
	quote.Status = OffState
	quote.OffReason = req.Reason
	res, _ := models.Show().Update(&quote)
	if res <= 0 {
		resp.FailJson(c, gin.H{}, -1, "下架失败")
		return
	}

	resp.SuccessJson(c, gin.H{})
}
