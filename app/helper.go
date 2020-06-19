// 框架函数库

package app

import "goweb/app/config"

// 取配置值，使用者需自己断言真实类型
func Get(prefix string, key string) interface{} {
	return config.Get(prefix, key)
}
