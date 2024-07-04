package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/anojaryal/Cancer-Cell-Detector/initializers"
	"github.com/anojaryal/Cancer-Cell-Detector/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY = os.Getenv("SECRET_KEY")

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
		HospitalID      int
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var existingUser models.User
	if initializers.DB.Where("username = ?", body.Username).First(&existingUser); existingUser.Username != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username already registered",
		})
		return
	}

	if initializers.DB.Where("email = ?", body.Email).First(&existingUser); existingUser.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already registered",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{
		Email:           body.Email,
		Username:        body.Username,
		FullName:        body.FullName,
		Address:         body.Address,
		BloodGroup:      body.BloodGroup,
		Gender:          body.Gender,
		ContactNo:       body.ContactNo,
		Password:        string(hash),
		IsVerified:      false,
		IsAdmin:         body.IsAdmin,
		IsHospitalAdmin: body.IsHospitalAdmin,
		HospitalID:      uint(body.HospitalID),
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	go func(email, token string) {
		err := sendVerificationEmail(email, token)
		if err != nil {
			fmt.Println("Failed to send verification email:", err)
		}
	}(user.Email, tokenString)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully. Please check your email for verification instructions.",
	})
}
