package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kljyrx/servicedesk/controllers"
	"github.com/kljyrx/servicedesk/helper"
	"github.com/kljyrx/servicedesk/models"
)

func main() {
	r := gin.Default()
	db := models.InitDB()
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
			helper.LogError(err.Error())
		}
	}(db)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/login", controllers.UserControl.Login)
	r.POST("/creatAdmin", controllers.UserControl.SaveUser)
	r.POST("/creatMachine", controllers.MachineControl.SaveMachine)
	r.POST("/listMachines", controllers.MachineControl.ListMachines)
	r.POST("/getMachineStatus", controllers.MachineControl.GetMachineStatus)
	err := r.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080
}
