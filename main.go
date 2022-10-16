package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kljyrx/servicedesk/controllers"
	"github.com/kljyrx/servicedesk/models"
)

func main() {
	r := gin.Default()
	db := models.InitDB()
	defer db.Close()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/login", controllers.UserContro.Login)
	r.POST("/creatAdmin", controllers.UserContro.SaveUser)
	r.Run() // listen and serve on 0.0.0.0:8080
}
