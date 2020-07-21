// 业务函数库

package helper

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"goweb/pkg/hot"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// 递归读取目录及所有子目录
// except 用于去除 path 中的某些字符，比如，.makrdown/PHP/Laravel 想要去除 .markdown
func RecursiveGetDirList(path string, except string) ([]map[string]interface{}, error) {
	// 读目录
	list, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, errors.New("文件系统目录读取失败" + path)
	}

	// 遍历及递归
	dirs := make([]map[string]interface{}, 0)
	for _, v := range list {
		if v.IsDir() && v.Name() != ".git" {
			name := v.Name()
			ut := v.ModTime().Format(hot.GetTimeCommonFormat())
			exportedPath := path + name + "/"
			if except != "" {
				// 目前只处理一次，主要为了将博客系统的根目录隐藏，有其他需要，可以再做扩展
				exportedPath = FormatPath(exportedPath, except)
			}
			// 遍历结构内无需再处理 error，该目录肯定存在
			son, _ := RecursiveGetDirList(path+name+"/", except)
			dirs = append(dirs, map[string]interface{}{
				"name": name,
				"ut":   ut,
				"son":  son,
				"path": exportedPath,
			})
		}
	}

	return dirs, nil
}

// 获取根目录
func GetBlogRoot() (string, error) {
	rootDir, _ := hot.GetConfig("blog.root").(string)
	if rootDir == "" {
		t := "博客系统配置缺失，应配置 blog.root，另需确保配置环境匹配"
		logrus.Errorln(t)
		return "", errors.New(t)
	}

	return rootDir, nil
}

// 获取格式 path
func FormatPath(path, except string) string {
	s := strings.Replace(path, except, "", 1)
	if len(s) == 0 {
		return ""
	}
	return string(s[0 : len(s)-1])
}

// 生成 jwt token
func JwtToken(s string) (string, int) {
	expired, _ := strconv.Atoi(os.Getenv("EXPIRES_AT"))
	// 元数据，对 token 的修饰，如过期时间、发行人
	claims := &jwt.StandardClaims{
		ExpiresAt: int64(expired),
		Issuer:    os.Getenv("ISSUER"),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte(s))
	return t, expired
}

// 验证 jwt token
func JwtCheck(s string) error {
	token, err := jwt.ParseWithClaims(s, &jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// since we only use the one private key to sign the tokens,
			// we also only use its public counter part to verify
			return verifyKey, nil
		})
}
