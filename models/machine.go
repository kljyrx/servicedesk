package models

import (
	"errors"
	"github.com/kljyrx/servicedesk/helper"
)

type Machine struct {
	ID         int    `gorm:"column:id; primary_key" json:"id"`
	Name       string `gorm:"column:name" json:"name"`
	User       string `gorm:"column:user" json:"user"`
	Host       string `gorm:"column:host" json:"host"`
	Port       string `gorm:"column:port" json:"port"`
	PassWord   string `gorm:"column:password" json:"password"`
	OperatorId int    `gorm:"column:operator_id; " json:"operator_id"`
}

type Machines []Machine

// TableName 自定义表名
func (m *Machine) TableName() string {
	return "machine"
}

func (m *Machine) FindMachineById() error {
	if err := db.First(m, "operator_id=? and id=?", m.OperatorId, m.ID).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("db查询失败")
	}
	return nil
}

func (m *Machine) AddMachine() error {
	if err := db.Create(m).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("添加机器失败")
	}
	return nil
}

func (m *Machines) ListMachines(operatorId int) error {
	if err := db.Find(m, "operator_id=?", operatorId).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("db查询失败")
	}
	return nil
}

func (m *Machines) ListMachinesByIds(operatorId int, ids []int) error {
	if len(ids) == 0 {
		return errors.New("ids值为空")
	}
	if err := db.Find(m, "operator_id=? and id in (?)", operatorId, ids).Error; err != nil {
		helper.LogError(err.Error())
		return errors.New("db查询失败")
	}
	return nil
}
