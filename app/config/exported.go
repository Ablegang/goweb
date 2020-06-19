// 使用说明：
// 1、	go 常驻内存且是编译型语言，每次在源码里修改配置后，需重新编译运行才会生效
// 2、	当前文件是 config 包唯一的出口文件，对外只开放 Config 函数，不开放 set 系列的函数
//	  	以保证所有 goroutine 得到的都是一份统一的配置
// 3、	若需要为特定的 goroutine 设置临时性的配置项，则在该 goroutine 内单独处理，copy 一份 c 的值，不要污染了 c 变量

package config

// c 是包全局变量，会在包被 import 时就初始化
var c = make(map[string]map[string]interface{})

// 取配置
func Get(prefix string, key string) interface{} {
	return c[prefix][key]
}
