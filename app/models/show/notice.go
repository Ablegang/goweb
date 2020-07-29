package show

import "time"

type Notice struct {
	Id        int64
	Key       string `xorm:"char(7) notnull comment('标识')"`
	Per       int64 `xorm:"notnull comment('幅度')"`
	CreatedAt time.Time
}

func (n *Notice) TableName() string {
	return "notices"
}
