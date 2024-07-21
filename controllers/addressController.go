package controllers

import (
	"net/http"
	"strconv"

	"github.com/anojaryal/Cancer-Cell-Detector/initializers"
	"github.com/anojaryal/Cancer-Cell-Detector/middleware"
	"github.com/anojaryal/Cancer-Cell-Detector/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AddPatientAddress
func AddPatientAddress(c *gin.Context) {
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

	currentUser, exists := c.Get(middleware.CurrentUser)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user := currentUser.(*models.User)

	if !user.IsAdmin && user.HospitalID != uint(hospitalID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var hospital models.Hospital
	if err := initializers.DB.First(&hospital, hospitalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hospital not found"})
		return
	}

	var patient models.Patient
	if err := initializers.DB.Where("hospital_id = ? AND id = ?", hospitalID, patientID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	var address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}
	if err := c.BindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	newAddress := models.Address{
		Street:    address.Street,
		City:      address.City,
		PatientID: patient.ID,
	}

	if err := initializers.DB.Create(&newAddress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Address added successfully",
		"address": newAddress,
	})
}

// GetPatientAddress
func GetPatientAddressByID(c *gin.Context) {
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

	addressIDStr := c.Param("address_id")
	addressID, err := uuid.Parse(addressIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	currentUser, exists := c.Get(middleware.CurrentUser)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user := currentUser.(*models.User)

	if !user.IsAdmin && user.HospitalID != uint(hospitalID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var patient models.Patient
	if err := initializers.DB.Where("hospital_id = ? AND id = ?", hospitalID, patientID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	var address models.Address
	if err := initializers.DB.Where("patient_id = ? AND id = ?", patientID, addressID).First(&address).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"address": address})
}

// UpdateAddress
func UpdatePatientAddress(c *gin.Context) {
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

	addressIDStr := c.Param("address_id")
	addressID, err := uuid.Parse(addressIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	currentUser, exists := c.Get(middleware.CurrentUser)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user := currentUser.(*models.User)

	if !user.IsAdmin && user.HospitalID != uint(hospitalID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var patient models.Patient
	if err := initializers.DB.Where("hospital_id = ? AND id = ?", hospitalID, patientID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	var address models.Address
	if err := initializers.DB.Where("id = ? AND patient_id = ?", addressID, patientID).First(&address).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	var input struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	address.Street = input.Street
	address.City = input.City

	if err := initializers.DB.Save(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"address": address})
}

// DeleteAddress
func DeleteAddress(c *gin.Context) {
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
	addressIDStr := c.Param("address_id")
	addressID, err := uuid.Parse(addressIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	currentUser, exists := c.Get(middleware.CurrentUser)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user := currentUser.(*models.User)

	if !user.IsAdmin && user.HospitalID != uint(hospitalID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var hospital models.Hospital
	if err := initializers.DB.First(&hospital, hospitalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hospital not found"})
		return
	}

	var patient models.Patient
	if err := initializers.DB.Where("hospital_id = ? AND id = ?", hospitalID, patientID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}
	var address models.Address
	if err := initializers.DB.Where("id = ? AND patient_id = ?", addressID, patientID).First(&address).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	if err := initializers.DB.Delete(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
}
