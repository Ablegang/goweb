package config

import "os"

func init() {
	c["app"] = map[string]interface{}{

		"port": os.Getenv("PORT"),
	}
}
