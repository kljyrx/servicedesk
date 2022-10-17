package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/kljyrx/servicedesk/helper"
	"github.com/kljyrx/servicedesk/models"
	"regexp"
	"strconv"
	"strings"
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
	var data []MachinesStatus
	for _, machine := range machines {
		passWord := helper.AesDecrypt(machine.PassWord)
		cli := helper.SSHCli{
			Addr: machine.Host + ":" + machine.Port,
			User: machine.User,
			Pwd:  passWord,
		}
		// 建立连接对象
		conncet, _ := cli.Connect()
		ret1, _ := conncet.Run("free -m|grep Mem")
		//Mem:           1993        1367         100           6         525         461
		regexp1, err := regexp.Compile(`^Mem:\s*(\d*)\s*(\d*)`)
		if err != nil {
			c.JSON(400, Response{Message: err.Error()})
			return
		}
		mem := regexp1.FindStringSubmatch(ret1)
		conncet.Run("export TERM=xterm")
		ret2, _ := conncet.RunTerminal("top -b -n 1|sed -n 3p")
		//%Cpu(s):  6.6 us,  4.5 sy,  0.0 ni, 88.4 id,  0.5 wa,  0.0 hi,  0.0 si,  0.0 st
		cpu := ret2[10:strings.Index(ret2, " us")]
		conncet.Client.Close()
		var machinesStatus MachinesStatus
		machinesStatus.Mem = helper.Division(mem[2], mem[1]) * 100
		machinesStatus.Cpu, _ = strconv.ParseFloat(cpu, 64)
		data = append(data, machinesStatus)
	}
	c.JSON(200, ResponseMachineStatus{Response: Response{Message: "获取机器信息成功！"}, Data: data})
}
