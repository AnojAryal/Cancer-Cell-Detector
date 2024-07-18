package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/anojaryal/Cancer-Cell-Detector/initializers"
	"github.com/anojaryal/Cancer-Cell-Detector/models"
	"github.com/anojaryal/Cancer-Cell-Detector/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

var SECRET_KEY = os.Getenv("SECRET_KEY")

func UserCreate(c *gin.Context) {
	// Retrieve the current user from the context
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user := currentUser.(*models.User)

	// Check if the current user is an admin or hospital admin
	if !user.IsAdmin && !user.IsHospitalAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

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

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var existingUser models.User
	if err := initializers.DB.Where("username = ?", body.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already registered"})
		return
	}

	if err := initializers.DB.Where("email = ?", body.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	//Determine hospital_id based on the current user's role
	hospitalID := 0
	if user.IsAdmin {
		hospitalID = body.HospitalID
	} else {
		hospitalID = int(user.HospitalID)
	}

	//Set is_admin to False by default if the current user is not an admin
	body.IsAdmin = user.IsAdmin

	user = &models.User{
		Email:           body.Email,
		Username:        body.Username,
		FullName:        body.FullName,
		Address:         body.Address,
		BloodGroup:      body.BloodGroup,
		Gender:          body.Gender,
		ContactNo:       body.ContactNo,
		Password:        string(hashedPassword),
		IsVerified:      false,
		IsAdmin:         body.IsAdmin,
		IsHospitalAdmin: body.IsHospitalAdmin,
		HospitalID:      uint(hospitalID),
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	go func(email, token string) {
		err := utils.SendVerificationEmail(email, token)
		if err != nil {
			fmt.Println("Failed to send verification email:", err)
		}
	}(user.Email, tokenString)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully. Please check your email for verification instructions.",
	})
}

// GetAllUsers
func GetAllUsers(c *gin.Context) {
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var users []models.User
	var query *gorm.DB

	// Check if the current user is an admin
	if currentUser.(*models.User).IsAdmin {
		query = initializers.DB.Find(&users)
	} else if currentUser.(*models.User).IsHospitalAdmin {
		hospitalID := currentUser.(*models.User).HospitalID
		query = initializers.DB.Where("hospital_id = ?", hospitalID).Find(&users)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
		return
	}

	// Execute the query and handle errors
	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserByID
func GetUserByID(c *gin.Context) {

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Retrieve current user from context
	currentUserInterface, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Assert currentUserInterface to models.User
	currentUser, ok := currentUserInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current user"})
		return
	}

	// Retrieve user from database
	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Admin can access all user details
	if currentUser.IsAdmin {
		c.JSON(http.StatusOK, user)
		return
	}

	// Hospital admin can access only users of their hospital
	if currentUser.IsHospitalAdmin && currentUser.HospitalID != user.HospitalID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this user's details"})
		return
	}

	// Ordinary user can access only their own details
	if currentUser.ID == uint(userID) {
		c.JSON(http.StatusOK, user)
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission"})

}

// PatchUserByID
func PatchUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var update_user struct {
		Username        *string
		Email           *string
		FullName        *string
		Address         *string
		BloodGroup      *string
		Gender          *string
		ContactNo       *string
		IsHospitalAdmin *bool
		IsAdmin         *bool
		HospitalID      *uint
	}

	if err := c.BindJSON(&update_user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	var user models.User
	if result := initializers.DB.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Retrieve current user from context
	currentUserInterface, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Assert currentUserInterface to models.User
	currentUser, ok := currentUserInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current user"})
		return
	}

	// Only admins or the hospital admin can update a user
	if !(currentUser.IsAdmin || (currentUser.IsHospitalAdmin && currentUser.HospitalID == user.HospitalID)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
		return
	}

	// Update fields
	if update_user.Username != nil {
		user.Username = *update_user.Username
	}
	if update_user.Email != nil {
		user.Email = *update_user.Email
	}
	if update_user.FullName != nil {
		user.FullName = *update_user.FullName
	}
	if update_user.Address != nil {
		user.Address = *update_user.Address
	}
	if update_user.BloodGroup != nil {
		user.BloodGroup = *update_user.BloodGroup
	}
	if update_user.Gender != nil {
		user.Gender = *update_user.Gender
	}
	if update_user.ContactNo != nil {
		user.ContactNo = *update_user.ContactNo
	}

	// Admin-specific updates
	if currentUser.IsAdmin {
		if update_user.IsHospitalAdmin != nil {
			user.IsHospitalAdmin = *update_user.IsHospitalAdmin
		}

		if update_user.IsAdmin != nil {
			user.IsAdmin = *update_user.IsAdmin
		}

		if update_user.HospitalID != nil {
			var hospital models.Hospital
			if result := initializers.DB.First(&hospital, *update_user.HospitalID); result.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Hospital not found"})
				return
			}
			user.HospitalID = *update_user.HospitalID
		}
	}

	// Hospital admin-specific updates
	if currentUser.IsHospitalAdmin {
		if update_user.HospitalID != nil && *update_user.HospitalID != currentUser.HospitalID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Hospital admins cannot change the hospital ID"})
			return
		}

		if update_user.IsHospitalAdmin != nil {
			user.IsHospitalAdmin = *update_user.IsHospitalAdmin
		}

		if update_user.IsAdmin != nil && !*update_user.IsAdmin {
			user.IsAdmin = *update_user.IsAdmin
		}
	}

	if result := initializers.DB.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// delete user by id
func DeleteUserByID(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	// Get the current user from context
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Assert the current user to *models.User (pointer to models.User)
	currentUserObj, ok := currentUser.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current user"})
		return
	}

	// Fetch the user to delete by ID
	if err := initializers.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		}
		return
	}

	// Only admins or the hospital admin associated with the user's hospital can delete the user
	if !(currentUserObj.IsAdmin || (currentUserObj.IsHospitalAdmin && currentUserObj.HospitalID == user.HospitalID)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
		return
	}

	// Delete the user
	if err := initializers.DB.Delete(&user, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
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

// get current user
func GetCurrentUser(c *gin.Context) {
	currentUserInterface, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userPtr, ok := currentUserInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current user"})
		return
	}
	user := *userPtr
	c.JSON(http.StatusOK, gin.H{
		"message": "Current user is: " + user.Username,
	})
}
