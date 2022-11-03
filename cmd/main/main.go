package main

import (
	"github.com/gin-gonic/gin"
	"github.com/teten-nugraha/golang-gin-project/pkg/config"
	"gorm.io/gorm"
)

var (
	db *gorm.DB = config.SetupDBConnection()
)

func main() {

	defer config.CloseDBConnection(db)

	app := gin.Default()
	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Up and running",
		})
	})
	app.Run()
}
