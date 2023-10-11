package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	apicontroller "github.com/mingyuanc/GovTech-Technical/controller/api_controller"
	"github.com/mingyuanc/GovTech-Technical/utils"
)

// Middleware to extract and validate query param
func ExtractAndValidateQueryTeacherParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the teacher query from the request
		teacherParam := c.QueryArray("teacher")
		if len(teacherParam) == 0 {
			c.IndentedJSON(400, gin.H{
				"error": "Required teacher query parameter not provided",
			})
			c.Abort()
		}
		// Validates the email address
		for i, teacher := range teacherParam {
			if !utils.IsValidEmail(teacher) {
				c.IndentedJSON(400, gin.H{
					"error": fmt.Sprintf("Teacher parameter at index %d is an invalid email: %s", i, teacher),
				})
				c.Abort()
			}

		}

		// Store the extracted parameter in the context for later use
		c.Set("teachersParam", teacherParam)

		// Call the next middleware or handler
		c.Next()
	}
}

// Adds the API routes to the current router
func AddApiRoutes(router *gin.Engine, conn *apicontroller.Connection) {
	api := router.Group("/api")
	// Able to add auth middleware here
	api.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	api.POST("/register", conn.HandleRegister)
	api.GET("/commonstudents", ExtractAndValidateQueryTeacherParam(), conn.HandleCommonStu)
}
