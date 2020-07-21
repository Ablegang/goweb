// 框架函数库

package helper

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// 判断数组下标是否存在
func IssetArrayIndex(arr []interface{}, index int) bool {
	for i := range arr {
		if i == index {
			return true
		}
	}

	return false
}

// 生成 MD5
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Sha1 加密
func Sha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	l := fmt.Sprintf("%x", h.Sum(nil))
	return l
}