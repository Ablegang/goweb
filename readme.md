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