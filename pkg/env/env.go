package env

import "github.com/joho/godotenv"

func init() {
	// 环境变量
	loadEnv()
}

// 载入 .env
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}
