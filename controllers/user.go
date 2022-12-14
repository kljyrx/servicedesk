package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/kljyrx/servicedesk/helper"
	"github.com/kljyrx/servicedesk/models"
)

type UserCreateRequest struct {
	Name     string `json:"name"`     // 用户名
	PassWord string `json:"password"` // 密码
}

type LoginRequest struct {
	Name     string `json:"name"`     // 用户名
	PassWord string `json:"password"` // 密码
}

type UserController struct {
}

func (t *UserController) Login(c *gin.Context) {
	helper.LogInfo("登录开始")
	var loginRequest LoginRequest

	// 将前端穿过来的json数据绑定存储在这个实体类中，BindJSON()也能使用
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		helper.LogError(err.Error())
		return
	}
	var user models.User
	user.Name = loginRequest.Name
	if err := user.FindUserByName(); err != nil {
		c.JSON(400, err)
		helper.LogError(err.Error())
		return
	}
	if user.PassWord != helper.Md5(loginRequest.PassWord) {
		c.JSON(400, Response{Message: "密码错误"})
		return
	}
	user.Token = helper.Md5(helper.Rand())
	if err := user.UpdateUser(); err != nil {
		c.JSON(400, Response{Message: "token设置错误"})
		return
	}
	c.JSON(200, ResponseLogin{Response: Response{Message: "登录成功！"}, Token: user.Token})
}

func (t *UserController) SaveUser(c *gin.Context) {
	var userCreateRequest UserCreateRequest

	// 将前端穿过来的json数据绑定存储在这个实体类中，BindJSON()也能使用
	if err := c.ShouldBindJSON(&userCreateRequest); err != nil {
		helper.LogError(err.Error())
		return
	}
	var user models.User
	if num := user.CountUserByName(userCreateRequest.Name); num > 0 {
		c.JSON(400, Response{Message: "用户名重复"})
		return
	}

	// 调用业务层的方法
	user.Name = userCreateRequest.Name
	user.PassWord = helper.Md5(userCreateRequest.PassWord)
	if err := user.AddUser(); err != nil {
		c.JSON(400, err)
		helper.LogError(err.Error())
		return
	}
	c.JSON(200, Response{Message: "添加成功"})
}
