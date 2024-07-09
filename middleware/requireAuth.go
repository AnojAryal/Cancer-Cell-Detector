package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/anojaryal/Cancer-Cell-Detector/initializers"
	"github.com/anojaryal/Cancer-Cell-Detector/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const CurrentUser = "currentUser"

func RequireAuth(c *gin.Context) {

	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Check token expiration
	exp, ok := claims["exp"].(float64)
	if !ok || float64(time.Now().Unix()) > exp {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Find user by token subject
	sub, ok := claims["sub"].(float64)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := initializers.DB.First(&user, int(sub)).Error; err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set(CurrentUser, user)

	// Set role flags in context
	isAdmin, _ := claims["is_admin"].(bool)
	isHospitalAdmin, _ := claims["is_hospital_admin"].(bool)
	hospitalID, _ := claims["hospital_id"].(float64)

	c.Set("is_admin", isAdmin)
	c.Set("is_hospital_admin", isHospitalAdmin)
	c.Set("hospital_id", uint(hospitalID))

	c.Next()
}
