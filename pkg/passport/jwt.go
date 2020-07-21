package passport

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

// Jwt 对象
type Jwt struct {
	SignKey []byte
}

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
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	t, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := t.Claims.(*CustomClaims); ok && t.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}

	return "", TokenError
}
