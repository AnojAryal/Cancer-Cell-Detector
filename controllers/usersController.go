package controllers

import (
	"net/http"

	"github.com/anojaryal/Cancer-Cell-Detector/initializers"
	"github.com/anojaryal/Cancer-Cell-Detector/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var body struct {
		Email           string
		Username        string
		FullName        string
		Address         string
		BloodGroup      string
		Gender          string
		ContactNo       string
		Password        string
		IsVerified      bool
		IsAdmin         bool
		IsHospitalAdmin bool
		HospitalID		int
		
	}

	// Bind request body to struct
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create user instance
	user := models.User{
		Email:           body.Email,
		Username:        body.Username,
		FullName:        body.FullName,
		Address:         body.Address,
		BloodGroup:      body.BloodGroup,
		Gender:          body.Gender,
		ContactNo:       body.ContactNo,
		Password:        string(hash),
		IsVerified:      body.IsVerified,
		IsAdmin:         body.IsAdmin,
		IsHospitalAdmin: body.IsHospitalAdmin,
		HospitalID: uint(body.HospitalID),
	}

	// Save user to database
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Respond with success message
	c.JSON(http.StatusCreated, gin.H{
		"message" : "User Created Successfully",
	})
}
