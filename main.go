package main

import (
	"fmt"
	"goweb/pkg/house"
)

func main() {
	driver, err := house.New(house.HzGovResource)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	list, err1 := driver.GetList("金雅苑", "1")
	if err1 != nil {
		fmt.Println(err1.Error())
	}

	if list != nil {
		for k,v := range list.List {
			fmt.Println(k,v)
		}
	}
	//app.Start()
}
