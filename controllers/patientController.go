package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anojaryal/Cancer-Cell-Detector/initializers"
	"github.com/anojaryal/Cancer-Cell-Detector/middleware"
	"github.com/anojaryal/Cancer-Cell-Detector/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreatePatient(c *gin.Context) {
	var patientCreate struct {
		FirstName string
		LastName  string
		Email     string
		Phone     string
		BirthDate time.Time
	}

	currentUser, exists := c.Get(middleware.CurrentUser)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user := currentUser.(*models.User)

	hospitalIDStr := c.Param("hospital_id")
	hospitalID, err := strconv.Atoi(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid hospital ID",
		})
		return
	}

	// Check if the hospital exists
	var hospital models.Hospital
	if err := initializers.DB.First(&hospital, hospitalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Hospital not found",
		})
		return
	}

	if err := c.BindJSON(&patientCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var hospitalIDToUse uint

	if user.IsAdmin {
		// Admin can create a patient for any hospital
		hospitalIDToUse = uint(hospitalID)
	} else {
		// Non-admin user can only create a patient for their associated hospital
		hospitalIDToUse = user.HospitalID
		if hospitalIDToUse != uint(hospitalID) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "You do not have permission to create a patient for this hospital",
			})
			return
		}
	}

	// Create a patient instance
	patient := models.Patient{
		FirstName:  patientCreate.FirstName,
		LastName:   patientCreate.LastName,
		Email:      patientCreate.Email,
		Phone:      patientCreate.Phone,
		BirthDate:  patientCreate.BirthDate,
		HospitalID: hospitalIDToUse,
	}

	result := initializers.DB.Create(&patient)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create patient",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Patient created successfully",
		"patient": patient,
	})
}

// GetPatients
func GetPatients(c *gin.Context) {
	hospitalIDStr := c.Param("hospital_id")
	hospitalID, err := strconv.Atoi(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid hospital ID",
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

	currentUser, exists := c.Get(middleware.CurrentUser)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user := currentUser.(*models.User)

	if !user.IsAdmin && user.HospitalID != uint(hospitalID) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Permission denied",
		})
		return
	}

	var patients []models.Patient
	result := initializers.DB.Where("hospital_id = ?", hospitalID).Find(&patients)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch patients",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"patients": patients,
	})
}

// GetPatient by id
func GetPatientById(c *gin.Context) {
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

	// Check if the hospital exists
	var hospital models.Hospital
	if err := initializers.DB.First(&hospital, hospitalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Hospital not found",
		})
		return
	}

	currentUser, exists := c.Get(middleware.CurrentUser)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user := currentUser.(*models.User)

	if !user.IsAdmin && user.HospitalID != uint(hospitalID) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Permission denied",
		})
		return
	}

	// Fetch the patient for the hospital
	var patient models.Patient
	if err := initializers.DB.Where("hospital_id = ? AND id = ?", hospitalID, patientID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Patient not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"patient": patient,
	})
}

// UpdatePatientById
func UpdatePatientById(c *gin.Context) {
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

	currentUser, exists := c.Get(middleware.CurrentUser)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user := currentUser.(*models.User)

	if !user.IsAdmin && user.HospitalID != uint(hospitalID) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Permission denied",
		})
		return
	}

	var patientUpdate struct {
		FirstName string
		LastName  string
		Email     string
		Phone     string
		BirthDate time.Time
	}

	if err := c.BindJSON(&patientUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Fetch the patient for the hospital
	var patient models.Patient
	if err := initializers.DB.Where("hospital_id = ? AND id = ?", hospitalID, patientID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Patient not found",
		})
		return
	}

	// Update patient fields
	patient.FirstName = patientUpdate.FirstName
	patient.LastName = patientUpdate.LastName
	patient.Email = patientUpdate.Email
	patient.Phone = patientUpdate.Phone
	patient.BirthDate = patientUpdate.BirthDate

	if result := initializers.DB.Save(&patient); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update patient",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Patient updated successfully",
		"patient": patient,
	})
}

// DeletePatient
func DeletePatientById(c *gin.Context) {
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

	currentUser, exists := c.Get(middleware.CurrentUser)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user := currentUser.(*models.User)

	if !user.IsAdmin && user.HospitalID != uint(hospitalID) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Permission denied",
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

	if err := initializers.DB.Delete(&patient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete patient",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Patient deleted successfully",
	})
}
