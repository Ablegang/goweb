package show

import "time"

type Quote struct {
	Id           int64
	Name         string  `xorm:"char(6) notnull comment('名称')"`
	Number       string  `xorm:"char(6) notnull comment('代码')"`
	Key          string  `xorm:"char(7) notnull comment('标识')"`
	InitialPrice float64 `xorm:"notnull comment('初选价')"`
	TodayPrice   float64 `xorm:"notnull comment('今日价')"`
	AddReason    string  `xorm:"varchar(255) notnull comment('入选理由')"`
	OffReason    string  `xorm:"varchar(255) notnull comment('下架理由')"`
	Status       string  `xorm:"char(10) notnull comment('状态')"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (quotes *Quote) TableName() string {
	return "quotes"
}
