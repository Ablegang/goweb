package cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

func Start() {
	c := cron.New()
	_, _ = c.AddFunc("30 * * * *", func() { fmt.Println("Every hour on the half hour") })
	c.Start()
}
