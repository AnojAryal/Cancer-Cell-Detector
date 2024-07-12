package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anojaryal/Cancer-Cell-Detector/initializers"
	"github.com/anojaryal/Cancer-Cell-Detector/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreatePatient
func CreatePatient(c *gin.Context) {
	var patientCreate struct {
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		Phone     string    `json:"phone"`
		BirthDate time.Time `json:"birth_date"`
	}

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

	// Create a patient instance
	patient := models.Patient{
		FirstName:  patientCreate.FirstName,
		LastName:   patientCreate.LastName,
		Email:      patientCreate.Email,
		Phone:      patientCreate.Phone,
		BirthDate:  patientCreate.BirthDate,
		HospitalID: uint(hospitalID),
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

	// Check if the hospital exists
	var hospital models.Hospital
	if err := initializers.DB.First(&hospital, hospitalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Hospital not found",
		})
		return
	}

	// Fetch patients for the hospital
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
