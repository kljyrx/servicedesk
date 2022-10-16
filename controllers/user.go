package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/kljyrx/servicedesk/models"
	"log"
)

type UserCreateRequest struct {
	Name     string `gorm:"column:name" json:"name"`         // 用户名
	PassWord string `gorm:"column:password" json:"password"` // 密码
}

type LoginRequest struct {
	Name     string `gorm:"column:name" json:"name"`         // 用户名
	PassWord string `gorm:"column:password" json:"password"` // 密码
}

type UserController struct {
}

func (t *UserController) Login(c *gin.Context) {
	var loginRequest LoginRequest

	// 将前端穿过来的json数据绑定存储在这个实体类中，BindJSON()也能使用
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		log.Fatal("参数错误")
		return
	}
	var user models.User
	user.Name = loginRequest.Name
	if err := models.FindUserByName(&user); err != nil {
		c.JSON(400, err)
		return
	}
	if user.PassWord != Md5(loginRequest.PassWord) {
		c.JSON(400, Response{Message: "密码错误"})
		return
	}
	user.Token = Md5(Rand())
	models.UpdateUser(&user)
	c.JSON(200, ResponseLogin{Response: Response{Message: "登录成功！"}, Token: user.Token})
}

func (t *UserController) SaveUser(c *gin.Context) {
	var userCreateRequest UserCreateRequest

	// 将前端穿过来的json数据绑定存储在这个实体类中，BindJSON()也能使用
	if err := c.ShouldBindJSON(&userCreateRequest); err != nil {
		log.Fatal("参数错误")
		return
	}

	if num := models.CountUserByName(userCreateRequest.Name); num > 0 {
		c.JSON(400, Response{Message: "用户名重复"})
		return
	}

	// 调用业务层的方法
	var user models.User
	user.Name = userCreateRequest.Name
	user.PassWord = Md5(userCreateRequest.PassWord)
	if err := models.AddUser(&user); err != nil {
		c.JSON(400, err)
		return
	}
	c.JSON(200, Response{Message: "添加成功"})
}
