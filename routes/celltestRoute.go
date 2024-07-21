package routes

import (
	"github.com/anojaryal/Cancer-Cell-Detector/controllers"
	"github.com/anojaryal/Cancer-Cell-Detector/middleware"
	"github.com/gin-gonic/gin"
)

func CellTestRoutes(r *gin.Engine) {
	authRequired := r.Group("/")
	authRequired.Use(middleware.RequireAuth)

	authRequired.POST("/hospital/:hospital_id/patients/:patient_id/cell_tests", controllers.CreateCellTest)
	authRequired.GET("/hospital/:hospital_id/patients/:patient_id/cell_tests", controllers.GetCellTests)
	authRequired.PUT("/hospital/:hospital_id/patients/:patient_id/cell_tests/:cell_test_id", controllers.UpdateCellTest)
	authRequired.DELETE("/hospital/:hospital_id/patients/:patient_id/cell_tests/:cell_test_id", controllers.DeleteCellTest)
	authRequired.POST("/hospital/:hospital_id/patients/:patient_id/cell_tests/:cell_test_id/data_images", controllers.PostImageData)
	authRequired.GET("/hospital/:hospital_id/patients/:patient_id/cell_tests/:cell_test_id/data_images", controllers.GetImageData)
	authRequired.POST("/hospital/:hospital_id/patients/:patient_id/cell_tests/:cell_test_id/results", controllers.PostResult)
	authRequired.GET("/hospital/:hospital_id/patients/:patient_id/cell_tests/:cell_test_id/results", controllers.GetResult)
	authRequired.POST("/hospital/:hospital_id/patients/:patient_id/cell_tests/:cell_test_id/results/:result_id/result_images", controllers.PostResultImage)
	authRequired.GET("/hospital/:hospital_id/patients/:patient_id/cell_tests/:cell_test_id/results/:result_id/result_images", controllers.GetResultImage)

}
