package routes

import (
	"github.com/anojaryal/Cancer-Cell-Detector/controllers"
	"github.com/gin-gonic/gin"
)

func HospitalRoutes(r *gin.Engine) {

	r.POST("/hospitals", controllers.CreateHospital)
	r.GET("/hospitals", controllers.GetAllHospitals)
	r.GET("/hospitals/:id", controllers.GetHospitalById)
	r.PUT("/hospitals/:id", controllers.UpdateHospitalById)
	r.DELETE("/hospitals/:id", controllers.DeleteHospitalById)

}
