package main

import (
	"github.com/anojaryal/Cancer-Cell-Detector/controllers"
	"github.com/anojaryal/Cancer-Cell-Detector/initializers"
	"github.com/anojaryal/Cancer-Cell-Detector/middleware"
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
	r.GET("/verify/:token", controllers.VerifyUserEmail)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/hospitals", controllers.CreateHospital)
	r.GET("/hospitals", controllers.GetHospitals)
	r.GET("/hospitals/:id", controllers.GetHospitalById)
	r.PUT("/hospitals/:id", controllers.UpdateHospitalById)
	r.DELETE("/hospitals/:id", controllers.DeleteHospitalById)

	r.Run()
}
