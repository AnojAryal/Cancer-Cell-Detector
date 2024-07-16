package controllers

import (
	"net/http"
	"strconv"

	"github.com/anojaryal/Cancer-Cell-Detector/initializers"
	"github.com/anojaryal/Cancer-Cell-Detector/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// creating celltest of a patient
func CreateCellTest(c *gin.Context) {
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

	var hospital models.Hospital
	if err := initializers.DB.First(&hospital, hospitalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Hospital not found",
		})
		return
	}

	var patient models.Patient
	if err := initializers.DB.Where("hospital_id = ? AND id = ?", hospitalID, patientID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Patient not found",
		})
		return
	}

	var celltest struct {
		Title           string `json:"Title"`
		Description     string `json:"Description"`
		DetectionStatus string `json:"DetectionStatus"`
		PatientID       int    `json:"PatientID"`
	}
	if err := c.BindJSON(&celltest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	newCelltest := models.CellTest{
		Title:           celltest.Title,
		Description:     celltest.Description,
		DetectionStatus: celltest.DetectionStatus,
		PatientID:       patientID,
	}

	if err := initializers.DB.Create(&newCelltest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create celltest",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":  "celltest created successfully",
		"celltest": newCelltest,
	})
}

// Get Patient Celltest
func GetCellTests(c *gin.Context) {
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

	var hospital models.Hospital
	if err := initializers.DB.First(&hospital, hospitalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Hospital not found",
		})
		return
	}

	var patient models.Patient
	if err := initializers.DB.Where("hospital_id = ? AND id = ?", hospitalID, patientID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Patient not found",
		})
		return
	}

	var celltest models.CellTest
	if err := initializers.DB.Where("patient_id = ?", patientID).First(&celltest).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Celltest not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"celltest": celltest,
	})
}

// Update Celltest
func UpdateCelltest(c *gin.Context) {
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

	celltestIDStr := c.Param("celltest_id")
	celltestID, err := uuid.Parse(celltestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid celltest ID",
		})
		return
	}

	var hospital models.Hospital
	if err := initializers.DB.First(&hospital, hospitalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Hospital not found",
		})
		return
	}

	var patient models.Patient
	if err := initializers.DB.Where("hospital_id = ? AND id = ?", hospitalID, patientID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Patient not found",
		})
		return
	}

	var celltest models.CellTest
	if err := initializers.DB.Where("patient_id = ? AND id = ?", patientID, celltestID).First(&celltest).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Celltest not found",
		})
		return
	}

	var celltest_update struct {
		Title           string `json:"Title"`
		Description     string `json:"Description"`
		DetectionStatus string `json:"DetectionStatus"`
		PatientID       int    `json:"PatientID"`
	}
	if err := c.BindJSON(&celltest_update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	celltest.Title = celltest_update.Title
	celltest.Description = celltest_update.Description
	celltest.DetectionStatus = celltest_update.DetectionStatus

	if err := initializers.DB.Save(&celltest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update celltest",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "celltest updated successfully",
		"celltest": celltest,
	})
}
