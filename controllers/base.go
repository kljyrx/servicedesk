package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kljyrx/servicedesk/models"
)

var (
	// UserControl 所有的controller类声明都在这儿
	UserControl    = &UserController{}
	MachineControl = &MachineController{}
)

type Response struct {
	Message string
}

type ResponseLogin struct {
	Response
	Token string
}

type ResponseListMachines struct {
	Response
	Machines models.Machines
}

type MachinesStatus struct {
	Mem   string
	Cpu   string
	Error string
}

type ResponseMachineStatus struct {
	Response
	Data []MachinesStatus
}

func Auth(c *gin.Context) (*models.User, error) {
	token := c.GetHeader("token")
	if token == "" {
		return nil, errors.New("token为空")
	}
	var user models.User
	user.Token = token
	if err := user.FindUserByToken(); err != nil {
		c.JSON(400, Response{Message: err.Error()})
		return nil, errors.New("没有权限")
	}
	return &user, nil
}
