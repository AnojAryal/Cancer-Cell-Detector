package routes

import (
	"github.com/anojaryal/Cancer-Cell-Detector/controllers"
	"github.com/anojaryal/Cancer-Cell-Detector/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	authRequired := r.Group("/")
	authRequired.Use(middleware.RequireAuth)

	authRequired.POST("/create-user", controllers.UserCreate)
	authRequired.GET("/users", controllers.GetAllUsers)
	authRequired.GET("/users/:id", controllers.GetUserByID)
	authRequired.PATCH("/users/:id", controllers.PatchUserByID)
	authRequired.DELETE("/users/:id", controllers.DeleteUserByID)
	authRequired.GET("/validate", controllers.Validate)
	authRequired.GET("/current-user", controllers.GetCurrentUser)

	r.GET("/verify/:token", controllers.VerifyUserEmail)
	r.POST("/login", controllers.Login)
}
