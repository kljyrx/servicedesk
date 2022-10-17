package models

import (
	"errors"
	"github.com/kljyrx/servicedesk/helper"
)

type User struct {
	ID       int    `gorm:"column:id; primary_key" json:"id"`
	Name     string `gorm:"column:name" json:"name"`         // 用户名
	PassWord string `gorm:"column:password" json:"password"` // 密码
	Token    string `gorm:"column:token" json:"token"`       // 密码
}

// TableName 自定义表名
func (u *User) TableName() string {
	return "admin_user"
}

func (u *User) CountUserByName(name string) int {
	var count int
	db.Model(u).Where("name = ?", name).Count(&count)
	return count
}

func (u *User) FindUserByName() error {
	if err := db.First(u, "name = ?", u.Name).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("db查询失败")
	}
	return nil
}

func (u *User) FindUserByToken() error {
	if err := db.First(u, "token=?", u.Token).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("db查询失败")
	}
	return nil
}

func (u *User) AddUser() error {
	if err := db.Create(u).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("添加用户失败")
	}
	return nil
}

func (u *User) UpdateUser() error {
	if err := db.Save(u).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("修改用户失败")
	}
	return nil
}
