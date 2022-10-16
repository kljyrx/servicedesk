package models

import (
	"errors"
	"log"
)

type User struct {
	ID       int    `gorm:"column:id; parimary_key", json:"id"`
	Name     string `gorm:"column:name" json:"name"`         // 用户名
	PassWord string `gorm:"column:password" json:"password"` // 密码
	Token    string `gorm:"column:token" json:"token"`       // 密码
}

func CountUserByName(name string) int {
	var count int
	db.Model(&User{}).Where("name <> ?", name).Count(&count)
	return count
}

func FindUserByName(user *User) error {
	if err := db.First(user, "name=?", user.Name).Error; err != nil {
		log.Fatal(err)
		return errors.New("db查询失败")
	}
	return nil
}

func AddUser(user *User) error {
	if err := db.Create(user).Error; err != nil {
		log.Fatal(err)
		return errors.New("添加用户失败")
	}
	return nil
}

func UpdateUser(user *User) error {
	if err := db.Save(user).Error; err != nil {
		log.Fatal(err)
		return errors.New("修改用户失败")
	}
	return nil
}
