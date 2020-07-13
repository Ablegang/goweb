package show

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"goweb/app/models"
	"goweb/app/models/show"
	"goweb/pkg/helper"
	"goweb/pkg/request"
	resp "goweb/pkg/response"
	"time"
)

type LoginForm struct {
	UserName string `json:"name" form:"name" validate:"max=20,min=2"`
	Pwd      string `json:"pwd" form:"pwd" validate:"max=30,min=6"`
}

func Login(c *gin.Context) {
	// 入参
	req := &LoginForm{}
	if err := request.Bind(c, req); err != nil {
		resp.FailJson(c, gin.H{}, -1, err.Error())
		return
	}

	// 查询
	admin := show.Admin{}
	query := models.Show().Where("name = ?", req.UserName)
	query = query.Or("email = ?", req.UserName).Or("phone = ?", req.UserName)
	has, _ := query.Get(&admin)
	if !has {
		resp.FailJson(c, gin.H{}, -1, "无此账号")
		return
	}

	// 校验密码
	ok := helper.CheckPwd(req.Pwd+admin.Salt, admin.Pwd)
	if !ok {
		resp.FailJson(c, gin.H{}, -1, "密码错误")
		return
	}

	// 登录成功，生成 Token
	info, _ := json.Marshal(admin)
	token, expiredAt := helper.JwtToken(string(info))

	// 更新登录时间
	admin.LastLoginAt = time.Now()
	_, err := models.Show().ID(admin.Id).Cols("last_login_at").Update(admin)
	if err != nil {
		logrus.Errorln("数据库更新失败", err)
	}

	// 响应
	resp.SuccessJson(c, gin.H{
		"token":     token,
		"expiredAt": expiredAt,
	})
}
