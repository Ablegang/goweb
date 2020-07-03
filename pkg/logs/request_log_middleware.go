// 这是日志组件内的小功能组件
// 任何 gin 框架都可以直接使用此组件来实现记录 Request 和 Response 日志

package logs

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"goweb/pkg/response"
	"io/ioutil"
	"time"
)

// 此变量主要用于确认依赖关系
// 当前包读取 response 内容，是依赖于 response 包的响应方法的
var _ = response.Rely

// 记录请求及响应日志
func RequestAndResponseLog(writer *CustomFileWriter) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()

		logJson := make(map[string]interface{})

		// 元数据
		logJson["@timestamp"] = start.String()
		logJson["@duration"] = end.Sub(start).String()
		logJson["@handlers"] = c.HandlerNames()
		logJson["@host"] = c.Request.Host
		logJson["@path"] = c.Request.URL.Path
		logJson["@uri"] = c.Request.Host + c.Request.URL.Path

		// 请求
		requestBody, _ := ioutil.ReadAll(c.Request.Body)
		logJson["@request"] = map[string]interface{}{
			"header": c.Request.Header,
			"query":  c.Request.URL.Query(),
			"body":   string(requestBody),
		}

		// 响应
		responseBody, _ := c.Get("responseBody")
		code, _ := c.Get("responseCode")
		header, _ := c.Get("responseHeader")
		logJson["@response"] = map[string]interface{}{
			"body":   responseBody,
			"code":   code,
			"header": header,
		}

		t, _ := json.Marshal(logJson)

		_, _ = writer.Write(t)
		_, _ = writer.Write([]byte("\n"))
	}
}
