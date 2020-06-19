// 日志相关的配置

package config

import (
	"github.com/sirupsen/logrus"
	"goweb/pkg/logs/loghooks"
)

func init() {

	c["log"] = map[string]interface{}{

		// 最小可用日志级别
		"minLevel": logrus.TraceLevel,

		// 是否报告打日志的位置
		"reportCaller": true,

		// hooks
		"hooks": []logrus.Hook{
			loghooks.NewEmailNotify(),
			loghooks.NewFileWriter(),
		},
	}

}
