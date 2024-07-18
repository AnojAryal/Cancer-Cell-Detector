package routes

import (
	"github.com/anojaryal/Cancer-Cell-Detector/controllers"
	"github.com/anojaryal/Cancer-Cell-Detector/middleware"
	"github.com/gin-gonic/gin"
)

func HospitalRoutes(r *gin.Engine) {
	adminRoutes := r.Group("/")
	adminRoutes.Use(middleware.RequireAuth)
	adminRoutes.Use(middleware.RequireAdmin)

	adminRoutes.POST("/hospitals", controllers.CreateHospital)
	adminRoutes.GET("/hospitals", controllers.GetAllHospitals)
	adminRoutes.GET("/hospitals/:id", controllers.GetHospitalById)
	adminRoutes.PUT("/hospitals/:id", controllers.UpdateHospitalById)
	adminRoutes.DELETE("/hospitals/:id", controllers.DeleteHospitalById)
}
