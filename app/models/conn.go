package models

import (
	_ "github.com/go-sql-driver/mysql"
	"goweb/app/models/show"
	"goweb/pkg/hot"
	"goweb/pkg/logs"
	"io"
	"os"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"
)

var (
	DBs = make(map[string]*xorm.Engine)

	Tables = map[string][]interface{}{
		// show 的表
		"show": []interface{}{
			&show.Admin{},
		},
	}
)

func Load() {
	// 取数据库配置
	cgs, _ := hot.GetConfig("mysqlDBs").([]interface{})
	if len(cgs) == 0 {
		logs.Panicln("数据库配置有误")
	}

	// 创建数据库引擎
	for _, v := range cgs {
		dbs, ok := v.(map[interface{}]interface{})
		if !ok {
			logs.Panicln("数据库配置格式有误，请检查")
		}

		for k, db := range dbs {
			name, _ := k.(string)
			dbMap, _ := db.(map[interface{}]interface{})
			dsn, _ := dbMap["dsn"].(string)
			engine, err := xorm.NewEngine("mysql", dsn)
			if err != nil {
				logs.Panicln("数据库引擎创建失败", err, dsn, name)
			}
			err = engine.Ping()
			if err != nil {
				logs.Panicln("数据库引擎连接失败", err, dsn, name)
			}

			// SQL 日志 writer
			engine.SetLogger(log.NewSimpleLogger2(io.MultiWriter(&logs.CustomFileWriter{
				LogMode:          os.Getenv("SQL_LOG_MODE"),
				Dir:              "storage/" + os.Getenv("SQL_LOG_DIR") + name + "/",
				FileNameFormater: os.Getenv("SQL_LOG_FILEFORMATER"),
				Perm:             os.FileMode(0777),
				IsDingRobot:      false,
			}, os.Stdout), "GoWeb-Show", log.DEFAULT_LOG_FLAG))
			// 日志级别
			engine.Logger().SetLevel(log.LOG_DEBUG)
			// 设置表名和结构体名映射算法，字段映射算法默认同
			engine.SetMapper(names.GonicMapper{})

			if auto, _ := dbMap["autoSync"].(string); auto == "on" {
				// 结构体和引擎关联，Sync2 会自动同步结构体里多出来的字段
				// 如果自动同步，表数量多的话，项目启动会很慢
				for _, t := range Tables[name] {
					if err := engine.Sync2(t); err != nil {
						logs.Panicln("表和结构体关联失败", err)
					}
				}
			}

			DBs[name] = engine
		}
	}
}

// show 连接
func Show() *xorm.Engine {
	return DBs["show"]
}
