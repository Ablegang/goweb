package house

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	HzGovResource = iota // 杭州房管局数据源
)

// 响应结构体
type Data struct {
	CityArea     string  `json:"cqmc"`     // 所属城区
	HouseLicense string  `json:"fczsh"`    // 房产证书号
	HouseCode    string  `json:"fwtybh"`   // 房源核验统一编码
	SellManName  string  `json:"gplxrxm"`  // 挂牌人姓名
	HouseArea    float64 `json:"jzmj"`     // 房产大小
	UploadDate   string  `json:"scgpshsj"` // 挂牌时间
	SellPrice    int     `json:"wtcsjg"`   // 委托出售价格
}

// 列表结构体
type ListData struct {
	List []Data `json:"list"`
}

// 驱动接口
type Driver interface {
	GetList(keyword string, page string) (*ListData, error)
}

// 驱动
type HzGovDriver struct {
	Api         string
	SignId      string
	Hash        string
	Method      string
	ContentType string
	Threshold   string
	Salt        string
}

// 获取列表
// 早上10点之后访问，太早的话，会请求失败
func (h HzGovDriver) GetList(keyword string, page string) (*ListData, error) {
	value := url.Values{}
	value.Set("keywords", keyword)
	value.Set("page", page)
	value.Set("signid", h.SignId)
	value.Set("threshold", h.Threshold)
	value.Set("salt", h.Salt)
	value.Set("nonce", "0")
	value.Set("hash", h.Hash)

	resp, err := http.Post(h.Api, h.ContentType, strings.NewReader(value.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return nil, err1
	}

	var list = new(ListData)
	err = json.Unmarshal(body, list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// 工厂
func New(flag int) (Driver, error) {
	var d Driver
	switch flag {
	case HzGovResource:
		d = &HzGovDriver{
			Api:         "http://jjhygl.hzfc.gov.cn/webty/WebFyAction_getGpxxSelectList.jspx",
			SignId:      "ff80808166484c980166486b4e0b0023",
			Hash:        "0448c9b2298cc81d7e0b7a2ab77fcd9261f956537b0939664985b08a1bc4ce20",
			Method:      http.MethodPost,
			ContentType: "application/x-www-form-urlencoded; charset=UTF-8",
			Threshold:   "ff80808166484c980166486b4e0b0021",
			Salt:        "ff80808166484c980166486b4e0b0022",
		}
	default:
		return nil, errors.New("错误的驱动类型")
	}

	return d, nil
}
