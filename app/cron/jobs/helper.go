package jobs

import (
	"goweb/app/models"
	"goweb/app/models/show"
	"os"
)

func GetAtMobile() string {
	return os.Getenv("AT_MOBILE")
}

func GetRobotToken() string {
	return os.Getenv("QUOTES_DING_ACCESS_TOKEN")
}

// 取数据库的 keys
func GetKeys() (keys []string) {
	model := show.Quote{}
	_ = models.Show().Table(model.TableName()).Cols("key").Find(&keys)
	return
}
