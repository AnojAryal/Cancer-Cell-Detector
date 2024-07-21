package controllers

import (
	"net/http"
	"strconv"

	"github.com/anojaryal/Cancer-Cell-Detector/initializers"
	"github.com/anojaryal/Cancer-Cell-Detector/models"
	"github.com/anojaryal/Cancer-Cell-Detector/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PostResultImage
func PostResultImage(c *gin.Context) {

	hospitalIDStr := c.Param("hospital_id")
	hospitalID, err := strconv.Atoi(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid hospital ID",
		})
		return
	}

	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid patient ID",
		})
		return
	}

	celltestIDStr := c.Param("cell_test_id")
	celltestID, err := uuid.Parse(celltestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid cell test ID",
		})
		return
	}

	resultIDStr := c.Param("result_id")
	resultID, err := uuid.Parse(resultIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid result ID",
		})
		return
	}

	// Retrieve the current user from context
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "Unauthorized"})
		return
	}

	user, ok := currentUser.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "Unauthorized"})
		return
	}

	// Check if the user has permission
	if !user.IsAdmin && user.HospitalID != uint(hospitalID) {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Permission denied"})
		return
	}

	var celltest models.CellTest
	if err := initializers.DB.Where("patient_id = ? AND id = ?", patientID, celltestID).First(&celltest).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cell test not found",
		})
		return
	}

	var result models.Result
	if err := initializers.DB.Where("cell_test_id = ? AND id = ?", celltestID, resultID).First(&result).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Result not found",
		})
		return
	}

	// Retrieve the file from the request
	file, err := c.FormFile("result_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid result image data",
		})
		return
	}

	// Define the upload directory
	uploadDir := "media/images/result_images"

	// Save the image file
	savedImagePath, err := utils.SaveImage(c, file, uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save result image",
		})
		return
	}

	// Save the result image data to the database
	resultImageData := models.ResultImage{
		ResultID: resultID,
		Image:    savedImagePath,
	}
	if err := initializers.DB.Create(&resultImageData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save result image data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Result image uploaded successfully",
		"data":    resultImageData,
	})
}

// GetResultImage
func GetResultImage(c *gin.Context) {
	hospitalIDStr := c.Param("hospital_id")
	hospitalID, err := strconv.Atoi(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid hospital ID",
		})
		return
	}

	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid patient ID",
		})
		return
	}

	celltestIDStr := c.Param("cell_test_id")
	celltestID, err := uuid.Parse(celltestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid cell test ID",
		})
		return
	}

	resultIDStr := c.Param("result_id")
	resultID, err := uuid.Parse(resultIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid result ID",
		})
		return
	}

	// Retrieve the current user from context
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "Unauthorized"})
		return
	}

	user, ok := currentUser.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "Unauthorized"})
		return
	}

	// Check if the user has permission
	if !user.IsAdmin && user.HospitalID != uint(hospitalID) {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Permission denied"})
		return
	}

	var celltest models.CellTest
	if err := initializers.DB.Where("patient_id = ? AND id = ?", patientID, celltestID).First(&celltest).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cell test not found",
		})
		return
	}

	var resultImageData models.ResultImage
	if err := initializers.DB.Where("result_id = ?", resultID).First(&resultImageData).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Result image not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resultImageData,
	})
}
