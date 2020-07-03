package response

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"goweb/pkg/snowflake"
	"strconv"
)

// 雪花算法生成分布式唯一 ID
func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()
		worker, err := snowflake.NewWorker(1)
		if err != nil {
			logrus.Errorln("Request Id 生成失败", err)
			return
		}
		requestId := strconv.FormatInt(worker.Next(),10)
		c.Header("RequestId", requestId)
		c.Set("RequestId", requestId)
	}
}
