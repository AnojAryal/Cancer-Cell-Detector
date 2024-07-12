package routes

import (
	"github.com/anojaryal/Cancer-Cell-Detector/controllers"
	"github.com/anojaryal/Cancer-Cell-Detector/middleware"
	"github.com/gin-gonic/gin"
)

func PasswordHandlerRoutes(r *gin.Engine) {
	r.PUT("/password-change", middleware.RequireAuth, controllers.ChangePassword)
	r.POST("/send-reset-email", controllers.SendResetEmail)
	r.POST("/reset-password", controllers.ResetPassword)
}
