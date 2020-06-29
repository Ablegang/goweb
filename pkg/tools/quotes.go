package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func main() {
	u := make(map[string]map[string]interface{})
	// 沪以 0 开头，深以 1 开头
	url := "http://api.money.126.net/data/feed/0600345,0000001,0000300,0600036,1002776"
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body[21:len(body)-2], &u)
	for _, v := range u {
		name, _ := v["name"].(string)
		percent, _ := v["percent"].(float64)
		percentStr := strconv.FormatFloat(percent*100, 'f', 6, 64)
		now := v["price"].(float64)
		nowStr := strconv.FormatFloat(now, 'f', 6, 64)
		fmt.Println(name, percentStr+"%", nowStr)
	}
	_ = resp.Body.Close()
}
