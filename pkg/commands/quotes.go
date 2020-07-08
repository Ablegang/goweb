package main

import (
	"fmt"
	"goweb/app/cron/quotes"
)

func main() {
	// 直接使用 make quotes ，并未载入 .env ，所以直接写死
	u := quotes.GetQuotes("0600345,0000001,0000300,0600036,1002776,0603101,0600110,1000672,0603799,1000032")
	for _, v := range u {
		name, _, percentStr, _, nowStr := quotes.FormatQuotesCoreData(v)
		fmt.Println(name, percentStr+"%", nowStr)
	}
}
