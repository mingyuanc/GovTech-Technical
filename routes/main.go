package routes

import (
	"github.com/gin-gonic/gin"
	apicontroller "github.com/mingyuanc/GovTech-Technical/controller/api_controller"
	"github.com/mingyuanc/GovTech-Technical/routes/api"
	"gorm.io/gorm"
)

// Run will start the server
func Run(db *gorm.DB) {
	router := gin.Default()
	conn := apicontroller.NewConnection(db)
	getRoutes(router, conn)
	router.Run(":8080")
}

// getRoutes will create our routes of our entire application
// this way every group of routes can be defined in their own file
// so this one won't be so messy
func getRoutes(router *gin.Engine, conn *apicontroller.Connection) {
	api.AddApiRoutes(router, conn)
}
