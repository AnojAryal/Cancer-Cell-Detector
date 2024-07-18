package routes

import (
	"github.com/anojaryal/Cancer-Cell-Detector/controllers"
	"github.com/anojaryal/Cancer-Cell-Detector/middleware"
	"github.com/gin-gonic/gin"
)

func HospitalRoutes(r *gin.Engine) {
	authRequired := r.Group("/")
	authRequired.Use(middleware.RequireAuth)

	authRequired.POST("/hospitals", controllers.CreateHospital)
	authRequired.GET("/hospitals", controllers.GetAllHospitals)
	authRequired.GET("/hospitals/:id", controllers.GetHospitalById)
	authRequired.PUT("/hospitals/:id", controllers.UpdateHospitalById)
	authRequired.DELETE("/hospitals/:id", controllers.DeleteHospitalById)

}
