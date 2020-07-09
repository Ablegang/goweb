package main

import (
	"fmt"
	"goweb/app/cron/quotes"
	_ "goweb/pkg/env"
	"os"
)

func main() {
	// 直接使用 make quotes ，并未载入 .env ，所以直接写死
	u := quotes.GetQuotes(os.Getenv("QUOTES"))
	for _, v := range u {
		name, _, percentStr, _, nowStr := quotes.FormatQuotesCoreData(v)
		fmt.Println(name, percentStr+"%", nowStr)
	}
}
