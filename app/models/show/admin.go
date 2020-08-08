package show

import (
	"goweb/pkg/helper"
)

type Admin struct {
	Id          int64
	Name        string `xorm:"index"`
	Email       string `xorm:"unique"`
	Phone       string `xorm:"unique"`
	Salt        string
	Pwd         string `xorm:"varchar(200)"`
	LastLoginAt helper.JsonTime
	CreatedAt   helper.JsonTime
	UpdatedAt   helper.JsonTime
}

func (admin *Admin) TableName() string {
	return "admins"
}
