package routes

import (
	"github.com/anojaryal/Cancer-Cell-Detector/controllers"
	"github.com/anojaryal/Cancer-Cell-Detector/middleware"
	"github.com/gin-gonic/gin"
)

func CellTestRoutes(r *gin.Engine) {

	r.POST("/hospital/:hospital_id/patients/:patient_id/cell_tests", middleware.RequireAuth, controllers.CreateCellTest)
	r.GET("/hospital/:hospital_id/patients/:patient_id/cell_tests", middleware.RequireAuth, controllers.GetCellTests)
	r.PUT("/hospital/:hospital_id/patients/:patient_id/cell_tests/:cell_test_id", middleware.RequireAuth, controllers.UpdateCellTest)
	r.DELETE("/hospital/:hospital_id/patients/:patient_id/cell_tests/:cell_test_id", middleware.RequireAuth, controllers.DeleteCellTest)
	r.POST("/hospital/:hospital_id/patients/:patient_id/cell_tests/:cell_test_id/data_images", controllers.PostImageData)
	r.GET("/hospital/:hospital_id/patients/:patient_id/cell_tests/:cell_test_id/data_images", controllers.GetImageData)
	r.POST("/hospital/:hospital_id/patients/:patient_id/cell_tests/:cell_test_id/results", controllers.PostResult)
	r.GET("/hospital/:hospital_id/patients/:patient_id/cell_tests/:cell_test_id/results", controllers.GetResult)
	r.POST("/hospital/:hospital_id/patients/:patient_id/cell_tests/:cell_test_id/results/:result_id/result_images", controllers.PostResultImage)
	r.GET("/hospital/:hospital_id/patients/:patient_id/cell_tests/:cell_test_id/results/:result_id/result_images", controllers.GetResultImage)

}
