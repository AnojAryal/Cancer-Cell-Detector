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

// Create Cell test
func CreateCellTest(c *gin.Context) {
	hospitalIDStr := c.Param("hospital_id")
	hospitalID, err := strconv.Atoi(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid hospital ID"})
		return
	}

	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid patient ID"})
		return
	}

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

	if !user.IsAdmin && user.HospitalID != uint(hospitalID) {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Permission denied"})
		return
	}

	var hospital models.Hospital
	if err := initializers.DB.First(&hospital, hospitalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Hospital not found"})
		return
	}

	var patient models.Patient
	if err := initializers.DB.Where("hospital_id = ? AND id = ?", hospitalID, patientID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Patient not found"})
		return
	}

	var cellTestInput struct {
		Title           string `json:"title" binding:"required"`
		Description     string `json:"description" binding:"required"`
		DetectionStatus string `json:"detection_status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&cellTestInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Failed to read body"})
		return
	}

	newCellTest := models.CellTest{
		Title:           cellTestInput.Title,
		Description:     cellTestInput.Description,
		DetectionStatus: cellTestInput.DetectionStatus,
		PatientID:       patientID,
	}

	if err := initializers.DB.Create(&newCellTest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create cell test"})
		return
	}

	c.JSON(http.StatusCreated, newCellTest)
}

// Get Patient Celltest
func GetCellTests(c *gin.Context) {
	hospitalIDStr := c.Param("hospital_id")
	hospitalID, err := strconv.Atoi(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid hospital ID"})
		return
	}

	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid patient ID"})
		return
	}

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

	if !user.IsAdmin && user.HospitalID != uint(hospitalID) {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Permission denied"})
		return
	}

	var patient models.Patient
	if err := initializers.DB.Where("hospital_id = ? AND id = ?", hospitalID, patientID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Patient not found"})
		return
	}

	var cellTests []models.CellTest
	if err := initializers.DB.Where("patient_id = ?", patientID).Find(&cellTests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to retrieve cell tests"})
		return
	}

	c.JSON(http.StatusOK, cellTests)
}

// UpdateCellTest
func UpdateCellTest(c *gin.Context) {
	hospitalIDStr := c.Param("hospital_id")
	hospitalID, err := strconv.Atoi(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid hospital ID"})
		return
	}

	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid patient ID"})
		return
	}

	cellTestIDStr := c.Param("cell_test_id")
	cellTestID, err := uuid.Parse(cellTestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid cell test ID"})
		return
	}

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

	if !user.IsAdmin && user.HospitalID != uint(hospitalID) {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Permission denied"})
		return
	}

	var cellTestUpdate struct {
		Title           string `json:"title" binding:"required"`
		Description     string `json:"description" binding:"required"`
		DetectionStatus string `json:"detection_status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&cellTestUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Failed to read body"})
		return
	}

	var cellTest models.CellTest
	if err := initializers.DB.Where("id = ? AND patient_id = ?", cellTestID, patientID).First(&cellTest).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Cell test not found"})
		return
	}

	cellTest.Title = cellTestUpdate.Title
	cellTest.Description = cellTestUpdate.Description
	cellTest.DetectionStatus = cellTestUpdate.DetectionStatus

	if err := initializers.DB.Save(&cellTest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to update cell test"})
		return
	}

	c.JSON(http.StatusOK, cellTest)
}

// Delete CellTest
func DeleteCellTest(c *gin.Context) {
	hospitalIDStr := c.Param("hospital_id")
	hospitalID, err := strconv.Atoi(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid hospital ID"})
		return
	}

	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid patient ID"})
		return
	}

	cellTestIDStr := c.Param("cell_test_id")
	cellTestID, err := uuid.Parse(cellTestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid cell test ID"})
		return
	}

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

	if !user.IsAdmin && user.HospitalID != uint(hospitalID) {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Permission denied"})
		return
	}

	if err := initializers.DB.Where("id = ? AND patient_id = ?", cellTestID, patientID).Delete(&models.CellTest{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to delete cell test"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "CellTest deleted successfully",
	})
}

// PostImageData
func PostImageData(c *gin.Context) {
	hospitalIDStr := c.Param("hospital_id")
	hospitalID, err := strconv.Atoi(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hospital ID"})
		return
	}

	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	cellTestIDStr := c.Param("cell_test_id")
	cellTestID, err := uuid.Parse(cellTestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cell test ID"})
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

	// Check if the hospital exists
	var hospital models.Hospital
	if err := initializers.DB.First(&hospital, hospitalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hospital not found"})
		return
	}

	// Check if the patient exists
	var patient models.Patient
	if err := initializers.DB.Where("hospital_id = ? AND id = ?", hospitalID, patientID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// Check if the cell test exists
	var cellTest models.CellTest
	if err := initializers.DB.Where("id = ? AND patient_id = ?", cellTestID, patientID).First(&cellTest).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cell test not found"})
		return
	}

	// Retrieve the file from the request
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image data"})
		return
	}

	// Define the upload directory
	uploadDir := "media/images/test_images"

	// Save the image file
	savedImagePath, err := utils.SaveImage(c, file, uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// Save the image data to the database
	imageData := models.CellTestImage{
		CellTestID: cellTestID,
		Image:      savedImagePath,
	}
	if err := initializers.DB.Create(&imageData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image data"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Image data uploaded successfully",
		"data":    imageData,
	})
}

// GetImageData
func GetImageData(c *gin.Context) {
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
	var imageData []models.CellTestImage
	if err := initializers.DB.Where("cell_test_id = ?", celltestID).Find(&imageData).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cell test image data not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    imageData,
	})
}

// PostResult
func PostResult(c *gin.Context) {
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

	var result struct {
		Description string `json:"Description"`
	}

	if err := c.BindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	new_result := models.Result{
		Description: result.Description,
		CellTestID:  celltestID,
	}
	if err := initializers.DB.Create(&new_result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Created",
		"data":    new_result,
	})
}

// Get result based on hospital ID, patient ID, and cell test ID
func GetResult(c *gin.Context) {
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

	var patient models.Patient
	if err := initializers.DB.Where("hospital_id = ? AND id = ?", hospitalID, patientID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Patient not found",
		})
		return
	}

	var result models.Result
	if err := initializers.DB.Where("cell_test_id = ?", celltestID).First(&result).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Result not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})

}
