package passport

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	resp "goweb/pkg/response"
	"os"
	"strconv"
	"time"
)

// 错误信息常量
var (
	TokenExpired = errors.New("令牌已超时")
	TokenError   = errors.New("令牌无效")
)

// 自定义 Claims
type CustomClaims struct {
	Content interface{}
	jwt.StandardClaims
}

// 创建自定义 Claims 实例
func NewCustomClaims(content interface{}) (*CustomClaims, int64) {
	expired, _ := strconv.Atoi(os.Getenv("EXPIRES_DUR"))
	c := &CustomClaims{
		Content: content,
	}
	c.ExpiresAt = time.Now().Add(time.Duration(expired) * time.Second).Unix()
	c.Issuer = os.Getenv("ISSUER")

	return c, int64(expired)
}

// Jwt 对象
type Jwt struct {
	SignKey []byte
}

// 创建 Jwt 实例
func NewJwt() *Jwt {
	return &Jwt{
		[]byte(os.Getenv("JWT_KEY")),
	}
}

// 创建 Token
func (j *Jwt) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SignKey)
}

// 刷新 Token
func (j *Jwt) RefreshToken(token string) (string, error) {
	// 解析旧 token
	t, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})
	if err != nil {
		return "", err
	}

	// 验证旧 token
	if claims, ok := t.Claims.(*CustomClaims); ok && t.Valid {
		jwt.TimeFunc = time.Now
		// 重新创建 token
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}

	return "", TokenError
}

// 解析 Token
func (j *Jwt) ParseToken(token string) (*CustomClaims, error) {
	// 解析 token
	t, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})

	// 解析错误处理
	if err != nil {
		if e, ok := err.(*jwt.ValidationError); ok {
			if e.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else {
				return nil, TokenError
			}
		}
	}

	// 过滤
	if t == nil {
		return nil, TokenError
	}

	// 返回解析成功的结果
	if claims, ok := t.Claims.(*CustomClaims); ok && t.Valid {
		return claims, nil
	}

	return nil, TokenError
}

// gin 中间件
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if len(token) <= 25 {
			resp.FailJson(c, "请先登录！")
			c.Abort()
			return
		}

		// 解析
		j := NewJwt()
		claims, err := j.ParseToken(token)
		if err != nil {
			resp.FailJson(c, err.Error())
			c.Abort()
			return
		}

		// 往 context 植入当前登录用户
		c.Set("AuthUser", claims)

		c.Next()
	}
}
