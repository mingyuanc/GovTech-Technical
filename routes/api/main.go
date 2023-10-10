package api

import (
	"github.com/gin-gonic/gin"
	apicontroller "github.com/mingyuanc/GovTech-Technical/controller/api_controller"
)

func AddApiRoutes(router *gin.Engine, conn *apicontroller.Connection) {
	api := router.Group("/api")
	api.GET("/ping", func(context *gin.Context) {
		context.JSON(200, gin.H{"message": "ok"})
	})
	api.POST("/register", conn.HandleRegister)
}
