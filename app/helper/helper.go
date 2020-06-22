// 框架函数库

package helper

import (
	"goweb/pkg/hot"
)

// 取配置值，使用者需自己断言真实类型
// 此方法依赖 logrus，所以必须在 logrus 注册之后使用
func Get(name string) interface{} {
	return hot.Get(name)
}