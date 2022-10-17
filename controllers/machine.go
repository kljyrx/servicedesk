package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kljyrx/servicedesk/helper"
	"github.com/kljyrx/servicedesk/models"
)

type MachineController struct {
}

type MachineStatusRequest struct {
	Ids []string `json:"ids"`
}

func (m *MachineController) SaveMachine(c *gin.Context) {
	var user *models.User
	var err error
	if user, err = Auth(c); err != nil {
		c.JSON(400, Response{Message: err.Error()})
		return
	}
	var machine models.Machine
	// 将前端穿过来的json数据绑定存储在这个实体类中，BindJSON()也能使用
	if err = c.ShouldBindJSON(&machine); err != nil {
		helper.LogError(err.Error())
		return
	}
	machine.OperatorId = user.ID
	machine.PassWord = helper.AesEncrypt(machine.PassWord)
	if err = machine.AddMachine(); err != nil {
		c.JSON(400, Response{Message: "添加机器失败" + err.Error()})
		return
	}
	c.JSON(200, Response{Message: "添加机器成功！"})
}

func (m *MachineController) ListMachines(c *gin.Context) {
	var user *models.User
	var err error
	if user, err = Auth(c); err != nil {
		c.JSON(400, Response{Message: err.Error()})
		return
	}
	var machines models.Machines
	if err = machines.ListMachines(user.ID); err != nil {
		c.JSON(400, Response{Message: err.Error()})
		return
	}
	for i, _ := range machines {
		machines[i].PassWord = "" //避免暴露密码给前端
	}
	c.JSON(200, ResponseListMachines{Response: Response{Message: "获取机器列表成功！"}, Machines: machines})
}

func (m *MachineController) GetMachineStatus(c *gin.Context) {
	var user *models.User
	var err error
	if user, err = Auth(c); err != nil {
		c.JSON(400, Response{Message: err.Error()})
		return
	}
	var machineStatusRequest MachineStatusRequest
	if err := c.ShouldBindJSON(&machineStatusRequest); err != nil {
		helper.LogError(err.Error())
		return
	}
	var machines models.Machines
	if err = machines.ListMachinesByIds(user.ID, machineStatusRequest.Ids); err != nil {
		c.JSON(400, Response{Message: err.Error()})
		return
	}
	for _, machine := range machines {
		passWord := helper.AesDecrypt(machine.PassWord)
		cli := helper.SSHCli{
			Addr: machine.Host + ":" + machine.Port,
			User: machine.User,
			Pwd:  passWord,
		}
		// 建立连接对象
		c, _ := cli.Connect()
		res, _ := c.Run("ls")
		res1, _ := c.Run("pwd")
		c.Client.Close()
		fmt.Println(res)
		fmt.Println(res1)
	}
	c.JSON(200, ResponseListMachines{Response: Response{Message: "获取机器列表成功！"}, Machines: machines})
}
