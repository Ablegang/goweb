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
	user := show.User{}
	query := models.Show().Where("name = ?", req.UserName)
	query = query.Or("email = ?", req.UserName).Or("phone = ?", req.UserName)
	has, _ := query.Get(&user)
	if !has {
		resp.FailJson(c, gin.H{}, -1, "无此账号")
		return
	}

	// 校验密码
	ok := passport.CheckPwd(req.Pwd, user.Salt, user.Pwd)
	if !ok {
		resp.FailJson(c, gin.H{}, -1, "密码错误")
		return
	}

	// 登录成功，生成 Token
	cl, expiredDur := passport.NewCustomClaims(user)
	token, _ := passport.NewJwt().CreateToken(*cl)

	// 更新登录时间
	user.LastLoginAt = helper.JsonTime(time.Now())
	_, err := models.Show().ID(user.Id).Cols("last_login_at").Update(user)
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
type AddUserForm struct {
	Name  string `json:"name" form:"name" validate:"max=30,min=2"`
	Email string `json:"email" form:"email" validate:"omitempty,email"`
	Phone string `json:"phone" form:"phone" validate:"omitempty,min=11,max=11"`
	Pwd   string `json:"pwd" form:"pwd" validate:"max=30,min=6"`
}

// 添加用户处理
func AddUser(c *gin.Context) {
	// 入参
	req := &AddUserForm{}
	if err := request.Bind(c, req); err != nil {
		resp.FailJson(c, gin.H{}, -1, err.Error())
		return
	}

	// 唯一性校验
	user := show.User{}
	has, _ := models.Show().Where("email = ?", req.Email).Or("phone = ?", req.Phone).Get(&user)
	if has {
		resp.FailJson(c, gin.H{}, -1, "已存在手机号或 email")
		return
	}

	// 插入数据
	salt := helper.Md5(strconv.Itoa(rand.Intn(10000)))
	_, err := models.Show().Insert(&show.User{
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Salt:        salt,
		Pwd:         passport.Pwd(req.Pwd, salt),
		LastLoginAt: helper.JsonTime(time.Now()),
		CreatedAt:   helper.JsonTime(time.Now()),
		UpdatedAt:   helper.JsonTime(time.Now()),
	})
	if err != nil {
		resp.FailJson(c, gin.H{}, -1, err.Error())
		return
	}

	// 响应
	resp.SuccessJson(c, gin.H{})
	return
}

// 用户信息
func UserInfo(c *gin.Context) {
	userI, _ := c.Get("AuthUser")
	user, _ := userI.(map[string]interface{})
	delete(user, "Salt")
	delete(user, "Pwd")

	query := models.Show().Where("user_id = ?", user["Id"])
	total, _ := query.Where("role_id = ?", show.ADMIN).Count(&show.UserRole{})
	user["isAdmin"] = total

	// TODO:按过期时间查询
	query = models.Show().Where("user_id = ?", user["Id"])
	total, _ = query.Where("role_id = ?", show.VIP).Count(&show.UserRole{})
	user["isVip"] = total

	if total > 0 {
		// TODO:查询过期时间
		user["vipExpired"] = helper.JsonTime(time.Now())
	}

	resp.SuccessJson(c, user)
	return
}
