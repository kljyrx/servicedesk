package controllers

import (
	"fmt"
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

type MachineUploadFileRequest struct {
	Ids           []string `json:"ids"`
	LocalFilePath string   `json:"localFilePath"`
	RemotePath    string   `json:"remotePath"`
}

type MachineDownloadFileRequest struct {
	Id            int    `json:"id"`
	LocalFilePath string `json:"localFilePath"`
	RemotePath    string `json:"remotePath"`
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
	for i := range machines {
		machines[i].PassWord = "" //避免暴露密码给前端
	}
	c.JSON(200, ResponseListMachines{Response: Response{Message: "获取机器列表成功！"}, Machines: machines})
}

func (m *MachineController) loginMachine(c *gin.Context, machine models.Machine) (*helper.SSHCli, error) {
	passWord := helper.AesDecrypt(machine.PassWord)
	cli := helper.SSHCli{
		Addr: machine.Host + ":" + machine.Port,
		User: machine.User,
		Pwd:  passWord,
	}
	// 建立连接对象
	connect, err := cli.Connect()
	if err != nil {
		c.JSON(400, Response{Message: err.Error()})
	}
	return connect, err
}

func (m *MachineController) sftpMachine(c *gin.Context, machine models.Machine) (*helper.SftpCli, error) {
	sshCli, err := m.loginMachine(c, machine)
	if err != nil {
		return nil, err
	}
	sftpCli := helper.SftpCli{SshClient: sshCli.Client}
	if err = sftpCli.Connect(); err != nil {
		c.JSON(400, Response{Message: err.Error()})
		return nil, err
	}
	return &sftpCli, nil
}

func (m *MachineController) GetMachineStatus(c *gin.Context) {
	user, err := Auth(c)
	if err != nil {
		return
	}
	var machineStatusRequest MachineStatusRequest
	if err := c.ShouldBindJSON(&machineStatusRequest); err != nil {
		helper.LogError(err.Error())
		return
	}
	var machines models.Machines
	if err := machines.ListMachinesByIds(user.ID, machineStatusRequest.Ids); err != nil {
		c.JSON(400, Response{Message: err.Error()})
		return
	}
	var data []MachinesStatus
	for _, machine := range machines {
		sshCli, err := m.loginMachine(c, machine)
		if err != nil {
			return
		}
		ret1, _ := sshCli.Run("free -m|grep Mem")
		//Mem:           1993        1367         100           6         525         461
		regexp1, err := regexp.Compile(`^Mem:\s*(\d*)\s*(\d*)`)
		if err != nil {
			c.JSON(400, Response{Message: err.Error()})
			return
		}
		mem := regexp1.FindStringSubmatch(ret1)
		ret2, _ := sshCli.RunTerminal("top -b -n 2 -d 3")
		//%Cpu(s):  6.6 us,  4.5 sy,  0.0 ni, 88.4 id,  0.5 wa,  0.0 hi,  0.0 si,  0.0 st
		regexp2, err := regexp.Compile(`(%Cpu\(s\): .*st)[\s\S]*(%Cpu\(s\): .*st)[\s\S]*`)
		if err != nil {
			c.JSON(400, Response{Message: err.Error()})
			return
		}
		cpus := regexp2.FindStringSubmatch(ret2)
		fmt.Println(cpus)
		cpu := cpus[2][10:strings.Index(cpus[2], " us")]
		_ = sshCli.Client.Close()
		var machinesStatus MachinesStatus
		machinesStatus.Mem = helper.Division(mem[2], mem[1]) * 100
		machinesStatus.Cpu, _ = strconv.ParseFloat(cpu, 64)
		data = append(data, machinesStatus)
	}
	c.JSON(200, ResponseMachineStatus{Response: Response{Message: "获取机器信息成功！"}, Data: data})
}

func (m *MachineController) UploadFile(c *gin.Context) {
	user, err := Auth(c)
	if err != nil {
		return
	}
	var mUR MachineUploadFileRequest
	if err := c.ShouldBindJSON(&mUR); err != nil {
		helper.LogError(err.Error())
		return
	}
	var machines models.Machines
	if err := machines.ListMachinesByIds(user.ID, mUR.Ids); err != nil {
		c.JSON(400, Response{Message: err.Error()})
		return
	}
	for _, machine := range machines {
		sftpCli, err := m.sftpMachine(c, machine)
		if err != nil {
			return
		}
		if err = sftpCli.UploadFile(mUR.LocalFilePath, mUR.RemotePath); err != nil {
			c.JSON(400, Response{Message: err.Error()})
			return
		}
	}
	c.JSON(200, Response{Message: "上传成功！"})
}

func (m *MachineController) DownloadFile(c *gin.Context) {
	user, err := Auth(c)
	if err != nil {
		return
	}
	var mDR MachineDownloadFileRequest
	if err := c.ShouldBindJSON(&mDR); err != nil {
		helper.LogError(err.Error())
		return
	}
	var machine models.Machine
	machine.ID = mDR.Id
	machine.OperatorId = user.ID
	if err := machine.FindMachineById(); err != nil {
		c.JSON(400, Response{Message: err.Error()})
		return
	}

	sftpCli, err := m.sftpMachine(c, machine)
	if err != nil {
		return
	}
	if err = sftpCli.DownloadFile(mDR.LocalFilePath, mDR.RemotePath); err != nil {
		c.JSON(400, Response{Message: err.Error()})
		return
	}

	c.JSON(200, Response{Message: "下载成功！"})
}
