package routes

import (
	"github.com/gin-gonic/gin"
	apicontroller "github.com/mingyuanc/GovTech-Technical/src/controller/api_controller"
	"github.com/mingyuanc/GovTech-Technical/src/routes/api"
	"gorm.io/gorm"
)

// Starts the router and listens for requests
func Run(db *gorm.DB) {
	router := gin.Default()
	conn := apicontroller.NewConnection(db)
	getRoutes(router, conn)
	router.Run(":8282")
}

// Returns a router for running the tests
func RunTest(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	conn := apicontroller.NewConnection(db)
	getRoutes(router, conn)
	return router
}

// getRoutes will create our routes of our entire application
// this way every group of routes can be defined in their own file
func getRoutes(router *gin.Engine, conn *apicontroller.Connection) {
	api.AddApiRoutes(router, conn)
}
