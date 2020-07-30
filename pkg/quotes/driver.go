package quotes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	WyResource = iota
)

// 标的数据结构体
type QuoteData struct {
	Name        string
	Percent     float64
	PercentStr  string
	Number      string
	NowPrice    float64
	NowPriceStr string
	Key         string
}

// 驱动接口
type Driver interface {
	// 设置标的标识列表
	SetKeys([]string)
	// 获取标的标识列表
	GetKeys() []string
	// 获取驱动名称
	GetName() string
	// 获取标的数据
	GetQuotes() ([]QuoteData, error)
	// 整理为 map 格式
	GetMap() (map[string]QuoteData, error)
}

// 工厂
func New(flag int) Driver {
	var d Driver
	switch flag {
	case WyResource:
		d = &WyDriver{
			Api:  "http://api.money.126.net/data/feed/",
			Name: "网易数据源",
			Keys: []string{"0000001"},
		}
	default:
		d = &WyDriver{
			Api:  "http://api.money.126.net/data/feed/",
			Name: "网易数据源",
			Keys: []string{"0000001"},
		}
	}

	return d
}

// 网易数据源驱动
type WyDriver struct {
	Api  string
	Name string
	Keys []string
}

func (d *WyDriver) SetKeys(keys []string) {
	d.Keys = keys
}

func (d *WyDriver) GetKeys() []string {
	return d.Keys
}

func (d *WyDriver) GetName() string {
	return d.Name
}

func (d *WyDriver) GetQuotes() ([]QuoteData, error) {
	// 请求数据
	url := d.Api + strings.Join(d.Keys, ",")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析数据
	body, _ := ioutil.ReadAll(resp.Body)
	u := make(map[string]map[string]interface{})
	_ = json.Unmarshal(body[21:len(body)-2], &u)

	// 组织数据
	var quotes []QuoteData
	for _, k := range d.Keys {
		if _, ok := u[k]["name"]; !ok {
			continue
		}
		quote := QuoteData{}
		quote.Name, _ = u[k]["name"].(string)
		quote.Percent, _ = u[k]["percent"].(float64)
		quote.PercentStr = strconv.FormatFloat(quote.Percent*100, 'f', 2, 64) + "%"
		quote.Number, _ = u[k]["symbol"].(string)
		quote.NowPrice, _ = u[k]["price"].(float64)
		quote.NowPriceStr = strconv.FormatFloat(quote.NowPrice, 'f', 2, 64)
		quote.Key = k
		quotes = append(quotes, quote)
	}

	return quotes, nil
}

func (d *WyDriver) GetMap() (map[string]QuoteData, error) {
	qList, err := d.GetQuotes()
	if err != nil {
		return nil, err
	}

	// 整理为 map 格式
	qMap := make(map[string]QuoteData)
	for _, q := range qList {
		qMap[q.Key] = q
	}

	return qMap, nil
}
