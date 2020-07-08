package app

import (
	"goweb/app/cron/quotes"
)

func startCron() {
	// 行情监控
	go quotes.ListenQuotesNotice()
	go quotes.ListenQuotesCommonPush()
	go quotes.NearCloseNotice()
	go quotes.NearOpenNotice()

	// 协程始终启动，除非 main 进程关闭
	select {}
}
