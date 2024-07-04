package main

import (
	"github.com/anojaryal/Cancer-Cell-Detector/controllers"
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
	r.POST("/signup", controllers.SignUp)
	r.POST("/login",controllers.Login)
	r.POST("/hospitals", controllers.CreateHospital)
	r.GET("/hospitals", controllers.GetHospitals)
	
	r.Run() 
}