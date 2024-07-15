package routes

import (
	"github.com/anojaryal/Cancer-Cell-Detector/controllers"
	"github.com/gin-gonic/gin"
)

func PatientRoutes(r *gin.Engine) {

	r.POST("/hospital/:hospital_id/patient", controllers.CreatePatient)
	r.GET("/hospital/:hospital_id/patients", controllers.GetPatients)
	r.GET("/hospital/:hospital_id/patient/:patient_id", controllers.GetPatientById)
	r.PUT("/hospital/:hospital_id/patient/:patient_id", controllers.UpdatePatientById)
	r.DELETE("/hospital/:hospital_id/patient/:patient_id", controllers.DeletePatientById)
	r.POST("/hospital/:hospital_id/patient/:patient_id/address", controllers.AddPatientAddress)
	r.GET("/hospital/:hospital_id/patient/:patient_id/address/:address_id", controllers.GetPatientAddressById)
	r.PUT("/hospital/:hospital_id/patient/:patient_id/address/:address_id", controllers.UpdateAddress)

}
