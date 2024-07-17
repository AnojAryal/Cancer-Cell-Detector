package routes

import (
	"github.com/anojaryal/Cancer-Cell-Detector/controllers"
	"github.com/gin-gonic/gin"
)

func CellTestRoutes(r *gin.Engine) {

	r.POST("/hospital/:hospital_id/patients/:patient_id/cell_tests", controllers.CreateCellTest)
	r.GET("/hospital/:hospital_id/patients/:patient_id/cell_tests", controllers.GetCellTests)
	r.PUT("/hospital/:hospital_id/patients/:patient_id/cell_tests/:celltest_id", controllers.UpdateCelltest)
	r.DELETE("/hospital/:hospital_id/patients/:patient_id/cell_tests/:celltest_id", controllers.DeleteCellTest)
	r.POST("/hospital/:hospital_id/patients/:patient_id/cell_tests/:celltest_id/data_images", controllers.PostImageData)
	r.GET("/hospital/:hospital_id/patients/:patient_id/cell_tests/:celltest_id/data_images", controllers.GetImageData)
	r.POST("/hospital/:hospital_id/patients/:patient_id/cell_tests/:celltest_id/results", controllers.PostResult)
	r.GET("/hospital/:hospital_id/patients/:patient_id/cell_tests/:celltest_id/results", controllers.GetResult)

}
