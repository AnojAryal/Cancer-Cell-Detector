package main

import (
	"github.com/anojaryal/Cancer-Cell-Detector/initializers"
	"github.com/anojaryal/Cancer-Cell-Detector/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	routes.UserRoutes(r)
	routes.HospitalRoutes(r)
	routes.PasswordHandlerRoutes(r)

	r.Run()
}
