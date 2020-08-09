package show

import "goweb/pkg/helper"

type UserRole struct {
	UserId    int64
	RoleId    int64
	CreatedAt helper.JsonTime
	UpdatedAt helper.JsonTime
}

func (relation *UserRole) TableName() string {
	return "user_roles"
}
