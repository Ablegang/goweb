package cron

import (
	"goweb/app/models"
	"goweb/app/models/show"
	"os"
)

var (
	AtMobile    = os.Getenv("AT_MOBILE")
	RobotToken  = os.Getenv("QUOTES_DING_ACCESS_TOKEN")
)

// 取数据库的 keys
func GetKeys() (keys []string) {
	model := show.Quote{}
	_ = models.Show().Table(model.TableName()).Cols("key").Find(&keys)
	return
}
