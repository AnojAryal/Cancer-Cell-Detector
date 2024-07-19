package routes

import (
	"github.com/anojaryal/Cancer-Cell-Detector/controllers"
	"github.com/anojaryal/Cancer-Cell-Detector/middleware"
	"github.com/gin-gonic/gin"
)

func PatientRoutes(r *gin.Engine) {

	authRequired := r.Group("/")
	authRequired.Use(middleware.RequireAuth)

	authRequired.POST("/hospital/:hospital_id/patients", controllers.CreatePatient)
	authRequired.GET("/hospital/:hospital_id/patients", controllers.GetPatients)
	authRequired.GET("/hospital/:hospital_id/patients/:patient_id", controllers.GetPatientById)
	authRequired.PUT("/hospital/:hospital_id/patients/:patient_id", controllers.UpdatePatientById)
	authRequired.DELETE("/hospital/:hospital_id/patients/:patient_id", controllers.DeletePatientById)
	authRequired.POST("/hospital/:hospital_id/patients/:patient_id/address", controllers.AddPatientAddress)
	authRequired.GET("/hospital/:hospital_id/patients/:patient_id/address/:address_id", controllers.GetPatientAddressByID)
	authRequired.PUT("/hospital/:hospital_id/patients/:patient_id/address/:address_id", controllers.UpdateAddress)
	authRequired.DELETE("/hospital/:hospital_id/patients/:patient_id/address/:address_id", controllers.DeleteAddress)

}
