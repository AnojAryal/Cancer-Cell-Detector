package routes

import (
	"github.com/anojaryal/Cancer-Cell-Detector/controllers"
	"github.com/gin-gonic/gin"
)

func CellTestRoutes(r *gin.Engine) {

	r.POST("/hospital/:hospital_id/patient/:patient_id/celltest", controllers.CreateCellTest)
	r.GET("/hospital/:hospital_id/patient/:patient_id/celltest", controllers.GetCellTests)

}
