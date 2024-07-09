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

func UserCreate(c *gin.Context) {
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

// VerifyUserEmail
func VerifyUserEmail(c *gin.Context) {
	token := c.Param("token")

	// Parse and validate the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	// Check if token is valid
	if !parsedToken.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	// Extract user ID from claims
	userID := int(parsedToken.Claims.(jwt.MapClaims)["user_id"].(float64))

	// Fetch user from database
	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if user is already verified
	if user.IsVerified {
		c.JSON(http.StatusOK, gin.H{"message": "Email is already verified"})
		return
	}

	// Mark user as verified and update in the database
	user.IsVerified = true
	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func GetCurrentUser(c *gin.Context) {
	// Retrieve current user from context
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user := currentUser.(models.User)

	c.JSON(http.StatusOK, gin.H{
		"message": "Current user is: " + user.Username,
	})
}
