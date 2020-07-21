package passport

import "goweb/pkg/helper"

// 检查密码格式
func CheckPwd(pwd, salt, sha1String string) bool {
	return Pwd(pwd, salt) == sha1String
}

// 生成密码
func Pwd(pwd, salt string) string {
	return helper.Sha1(pwd + salt)
}