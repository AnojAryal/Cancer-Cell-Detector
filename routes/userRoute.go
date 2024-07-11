package routes

import (
	"github.com/anojaryal/Cancer-Cell-Detector/controllers"
	"github.com/anojaryal/Cancer-Cell-Detector/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.POST("/create-user", middleware.RequireAuth, controllers.UserCreate)
	r.GET("/users", controllers.GetAllUsers)
	r.GET("/users/:id", controllers.GetUserByID)
	r.PATCH("/users/:id", controllers.PatchUserByID)
	r.DELETE("/users/:id", controllers.DeleteUserById)
	r.GET("/verify/:token", controllers.VerifyUserEmail)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/current-user", middleware.RequireAuth, controllers.GetCurrentUser)
}
