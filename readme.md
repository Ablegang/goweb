# 日志 Hook

包 ：pkg/logs/loghooks

日志 Hook 使用：
> 详情看源码

```go
package main
import (
    "os"
    "time"

    log "github.com/sirupsen/logrus"
    "goweb/pkg/logs/loghooks"
)

func init() {
    // 直接使用 log 包的函数，包内部会实例化一个 std 作为通用 logger
    // std 是指针类型，且是包变量，程序运行时就已经注册
    // 设置最小 log 级别
    log.SetLevel(log.TraceLevel)

    // 注册 Hooks...
    // 默认 FileWriter
    log.AddHook(loghooks.NewFileWriter())

    // 注册自定义 FileWriter（示例使用）
    log.AddHook(&loghooks.FileWriter{
        LogMode:          "single",
        Dir:              "storage/logs-test/",
        FileNameFormater: "logs.txt",
        EntryFormatter: &log.JSONFormatter{
            TimestampFormat: time.RFC3339Nano, // 含纳秒
        },
        HookLevels: log.AllLevels,
        Perm:       os.FileMode(0777),
    })
}
```

# 配置
当前配置分为两种，.env 配置和 yaml 配置
初次启动时（主进程）需要用到的配置，都写在 env 里，而处理请求时要用到的配置，则写在 yaml 里
yaml 支持热更新，只要修改 .version 的版本号即可
env 不支持热更新

> 关于 mysql 配置、redis 配置等变量，也完全可以写在 yaml 配置中，以此支持热更新，只要在 .gitignore 里管理好即可保证安全

# 用文件记录 gin 请求日志
> app/boot.go 的 init 函数
```go

    // 自定义 gin logger ，主要为了将日志数据按日输出到自定义的目录内
	writer := io.MultiWriter(gin.DefaultWriter, logs.NewCustomFileWriter()) // 可以看 logs/custom_file_writer.go
	r.Use(gin.LoggerWithWriter(writer))
```
# makefile
可参考 makefile 和 ./shell 下的文件
关于 git 的操作，应提前建好 git 项目，且不能为初始的 git 项目，必须有提交记录（远程仓库也一样）

# BLOG 设计
- 所有文章都使用 markdown 编写，md 文件存储，github 托管
- 文章内容单独放一个项目存储，由 spug 发布
- 规划文章结构，目录为分类，md 文件为文章，更新时间用 fileinfo
- 发布配置：
    - 排除文件：.git
    - 自定义全局变量：
        ```
        STORAGEPATH=/www/wwwroot/storage/goweb
        ENVPATH=/www/wwwroot/env/goweb
        BLOGPATH=/www/wwwroot/blog
        ```
    - 发布前钩子：
        ```
        # 注意这里的 env 和 blog 会作为单独 git 项目来管理，所以在发布之前，请先建立好这两个项目
        cd ${ENVPATH} && git pull origin master
        cd ${BLOGPATH} && git pull origin master
        ```
    - 发布后钩子：
        ```
          ln -s  ${STORAGEPATH} storage|| echo ""
          ln -s ${ENVPATH}/.env .env || echo ""
          ln -s ${ENVPATH}/prod.yaml prod.yaml || echo ""
          ln -s ${BLOGPATH} .markdown || echo ""
          go build
          
          # 停服
          PID=$(ps -e|grep 'goweb' |grep -v grep|awk '{printf $1}')
          echo $PID
          kill -9 $PID || echo ""
          
          nohup ./goweb &> start.log &
        ```
      
# storage
## logs
```
    custom 是开发者在程序中埋点打出的日志
    ginStd 是 gin 默认写出的请求日志
    ginErr 是 gin 遭遇 panic 的日志
    requests 是详细的请求及响应日志
```
# response

所有 response 相关的操作都封装在该包内，如果要改默认的返回格式，可以改 SuccessJson 或 FailJson 函数

另外，response 包重写了 gin 的 Recovery 中间件，实现记录到文件，同时所有的错误都统一抛出 404 ，而不是 500

使用 response.FailJson 时，一定要判断，相关报错信息是否需要返回给前端，如果不需要，则 msg 直接使用默认的 fail 即可

# 增加 RequestId 中间件
基于雪花算法
gin.Context.Get("RequestId") 即可在其他 Handlers 中获取

# 钉钉群机器人告警
- 在 logrus 注册了 hook，详情查看 pkg/logs/loghooks/ding_notify.go
- 为 gin 的 recover 上了钉钉群告警

# 框架规划
所有业务逻辑写在 app/Handlers 下
所有组件类功能，写在 pkg 下，如与业务无关，可以直接复制到其他项目中使用的，就可以定义为组件
后期如果接口多了，业务复杂了，则一个接口一个文件，相关表单校验、业务处理等所有内容，都写在这个文件里

# cron 计时任务
在 app/cron.go 里注册
