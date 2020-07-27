package main

import (
	"fmt"
	_ "goweb/pkg/env"
	"goweb/pkg/quotes"
	"os"
	"strings"
	"time"
)

func main() {
	driver := quotes.New(quotes.WyResource)

	keys := os.Getenv("QUOTES")
	driver.SetKeys(strings.Split(keys, ","))

	data, err := driver.GetQuotes()
	if err != nil {
		panic(err)
	}

	for _, v := range data {
		fmt.Printf("%s（%s）	%s 		%s \n", v.Name, v.Number, v.PercentStr, v.NowPriceStr)
	}

	fmt.Println(time.Now().String())
}
