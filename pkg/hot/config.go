// 此文件负责配置的热更新
// 基于根目录下的 .version

package hot

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

var oldVersion = ""
var config = make(map[interface{}]interface{})

// 配置更新
func fresh() {
	newVersion := getLastVersion()
	if oldVersion != newVersion {
		loadYaml()
		oldVersion = newVersion
	}
}

// 载入配置
func loadYaml() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "prod"
	}

	file := env + ".yaml"
	b, err := ioutil.ReadFile(file)
	if err != nil {
		logrus.Errorln("配置文件读取失败：", err)
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		logrus.Errorln("配置文件转码失败：", err)
	}
}

// 取配置项
func Get(name string) interface{} {
	// 检查更新
	fresh()

	// 支持 . 语法
	path := strings.Split(name, ".")
	data := config

	for key, value := range path {
		v, ok := data[value]
		if !ok {
			break
		}

		if (key + 1) == len(path) {
			return v
		}
		if reflect.TypeOf(v).String() == "map[interface {}]interface {}" {
			data = v.(map[interface{}]interface{})
		}
	}

	return nil
}

// 获取最新配置版本号
func getLastVersion() string {
	data, err := godotenv.Read(".version")
	if err != nil {
		logrus.Errorln("无法读取 .version 文件", err)
	}

	return data["CONFIG_VERSION"]
}
