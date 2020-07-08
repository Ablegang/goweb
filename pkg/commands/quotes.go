package main

import (
	"fmt"
	"goweb/app/cron"
	"strconv"
)

func main() {
	u := cron.GetQuotes("0600345,0000001,0000300,0600036,1002776,0603101,0600110,1000672,0603799,1000032")
	for _, v := range u {
		name, _ := v["name"].(string)
		percent, _ := v["percent"].(float64)
		percentStr := strconv.FormatFloat(percent*100, 'f', 6, 64)
		now := v["price"].(float64)
		nowStr := strconv.FormatFloat(now, 'f', 6, 64)
		fmt.Println(name, percentStr+"%", nowStr)
	}
}
