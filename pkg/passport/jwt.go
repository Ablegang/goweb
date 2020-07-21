package passport

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
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

