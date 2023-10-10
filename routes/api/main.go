package api

import (
	"github.com/gin-gonic/gin"
	apicontroller "github.com/mingyuanc/GovTech-Technical/controller/api_controller"
)

// Adds the API routes to the current router
func AddApiRoutes(router *gin.Engine, conn *apicontroller.Connection) {
	api := router.Group("/api")
	// Able to add auth middleware here
	api.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	api.POST("/register", conn.HandleRegister)
}
