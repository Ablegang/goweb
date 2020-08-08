package show

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"goweb/app/models"
	"goweb/app/models/show"
	"goweb/pkg/helper"
	"goweb/pkg/passport"
	"goweb/pkg/request"
	resp "goweb/pkg/response"
	"math/rand"
	"strconv"
	"time"
)

// 登录表单
type LoginForm struct {
	UserName string `json:"name" form:"name" validate:"max=30,min=2"`
	Pwd      string `json:"pwd" form:"pwd" validate:"max=30,min=6"`
}

// 登录处理
func Login(c *gin.Context) {

	// 预防爆破
	time.Sleep(3 * time.Second)

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
	ok := passport.CheckPwd(req.Pwd, admin.Salt, admin.Pwd)
	if !ok {
		resp.FailJson(c, gin.H{}, -1, "密码错误")
		return
	}

	// 登录成功，生成 Token
	cl, expiredDur := passport.NewCustomClaims(admin)
	token, _ := passport.NewJwt().CreateToken(*cl)

	// 更新登录时间
	admin.LastLoginAt = time.Now()
	_, err := models.Show().ID(admin.Id).Cols("last_login_at").Update(admin)
	if err != nil {
		logrus.Errorln("数据库更新失败", err)
	}

	// 响应
	resp.SuccessJson(c, gin.H{
		"token":      token,
		"expiredDur": expiredDur,
	})
}

// 添加用户表单
type AddAdminForm struct {
	Name  string `json:"name" form:"name" validate:"max=30,min=2"`
	Email string `json:"email" form:"email" validate:"omitempty,email"`
	Phone string `json:"phone" form:"phone" validate:"omitempty,min=11,max=11"`
	Pwd   string `json:"pwd" form:"pwd" validate:"max=30,min=6"`
}

// 添加用户处理
func AddAdmin(c *gin.Context) {
	// 入参
	req := &AddAdminForm{}
	if err := request.Bind(c, req); err != nil {
		resp.FailJson(c, gin.H{}, -1, err.Error())
		return
	}

	// 唯一性校验
	admin := show.Admin{}
	has, _ := models.Show().Where("email = ?", req.Email).Or("phone = ?", req.Phone).Get(&admin)
	if has {
		resp.FailJson(c, gin.H{}, -1, "已存在手机号或 email")
		return
	}

	// 插入数据
	salt := helper.Md5(strconv.Itoa(rand.Intn(10000)))
	_, err := models.Show().Insert(&show.Admin{
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Salt:        salt,
		Pwd:         passport.Pwd(req.Pwd, salt),
		LastLoginAt: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		resp.FailJson(c, gin.H{}, -1, err.Error())
		return
	}

	// 响应
	resp.SuccessJson(c, gin.H{})
	return
}
