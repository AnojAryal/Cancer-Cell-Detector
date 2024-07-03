package main

import (
	"github.com/anojaryal/Cancer-Cell-Detector/initializers"
	"github.com/gin-gonic/gin"
)


func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()	
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello! this is the strarting point of the app",
		})
	})
	r.Run()
}