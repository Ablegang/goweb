// 框架函数库

package app

import "goweb/app/config"

// 取配置值，使用者需自己断言真实类型
// 为了规范，一般业务代码里都使用 app 包的 Get 函数，即此函数，而不直接使用 config.Get
func Get(prefix string, key string) interface{} {
	return config.Get(prefix, key)
}
