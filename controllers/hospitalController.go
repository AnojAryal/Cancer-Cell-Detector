package controllers

import (
	"net/http"

	"github.com/anojaryal/Cancer-Cell-Detector/initializers"
	"github.com/anojaryal/Cancer-Cell-Detector/models"
	"github.com/gin-gonic/gin"
)

func CreateHospital(c *gin.Context) {
	// Define a struct for request body
	var body struct {
		Name	string 
		Address	string 
		Phone	string    
		Email	string 
	
	}

	// Bind request body to struct
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Create hospital instance
	hospital := models.Hospital{
		Name	:	body.Name,
		Address	:	body.Address,
		Phone	:	body.Phone,
		Email	:	body.Email,
	}

	// Save hospital to database
	result := initializers.DB.Create(&hospital)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create hospital",
		})
		return
	}

	// Respond with success message
	c.JSON(http.StatusCreated, gin.H{
		"message": "Hospital created successfully",
	})
}

//GET request to fetch all hospitals
func GetHospitals(c *gin.Context) {

	type Hospital struct {
		ID      uint   
		Name    string 
		Address string 
		Phone   string 
		Email   string
	}
	
	var hospitals []Hospital

	// Fetch hospitals from database
	result := initializers.DB.Find(&hospitals)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch hospitals",
		})
		return
	}

	// Respond with fetched hospitals
	c.JSON(http.StatusOK, hospitals)
}