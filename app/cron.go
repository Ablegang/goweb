package app

import "goweb/app/cron"

func startCron() {
	// 行情监控
	go cron.ListenQuotesNotice()
	go cron.ListenQuotesCommonPush()

	// 协程始终启动，除非 main 进程关闭
	select {}
}
