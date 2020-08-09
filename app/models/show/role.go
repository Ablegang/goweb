package show

import "goweb/pkg/helper"

type Role struct {
	Id        int64
	Name      string
	Intro     string
	CreatedAt helper.JsonTime
	UpdatedAt helper.JsonTime
}

func (role *Role) TableName() string {
	return "roles"
}

const (
	ADMIN = 1
	VIP   = 2
)
