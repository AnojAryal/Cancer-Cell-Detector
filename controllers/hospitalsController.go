package controllers

import (
	"net/http"

	"github.com/anojaryal/Cancer-Cell-Detector/initializers"
	"github.com/anojaryal/Cancer-Cell-Detector/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateHospital
func CreateHospital(c *gin.Context) {

	var body struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Phone   string `json:"phone"`
		Email   string `json:"email"`
	}

	// Bind request body to struct
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	//hospital instance
	hospital := models.Hospital{
		Name:    body.Name,
		Address: body.Address,
		Phone:   body.Phone,
		Email:   body.Email,
	}

	result := initializers.DB.Create(&hospital)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create hospital",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Hospital created successfully",
	})
}

// get all hospitals
func GetAllHospitals(c *gin.Context) {
	var hospitals []models.Hospital
	if result := initializers.DB.Preload("Users").Preload("Patients").Find(&hospitals); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve hospitals"})
		return
	}
	c.JSON(http.StatusOK, hospitals)
}

// GetHospitalById
func GetHospitalById(c *gin.Context) {
	var hospital models.Hospital

	// Extract the hospital ID from the URL parameters
	id := c.Param("id")

	// Fetch the hospital from the database by ID
	result := initializers.DB.First(&hospital, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Hospital not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hospital"})
		}
		return
	}

	c.JSON(http.StatusOK, hospital)
}

// update hospital by id
func UpdateHospitalById(c *gin.Context) {
	var hospital models.Hospital

	id := c.Param("id")

	if err := initializers.DB.First(&hospital, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Hospital not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hospital"})
		}
		return
	}
	if err := c.ShouldBindJSON(&hospital); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := initializers.DB.Save(&hospital).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hospital"})
		return
	}

	c.JSON(http.StatusOK, hospital)
}

func DeleteHospitalById(c *gin.Context) {
	var hospital models.Hospital

	id := c.Param("id")

	if err := initializers.DB.First(&hospital, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Hospital doesn't exist"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hospital"})
		}
		return
	}

	if err := initializers.DB.Delete(&hospital, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete hospital"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hospital deleted successfully"})
}
