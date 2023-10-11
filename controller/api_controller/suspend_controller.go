package apicontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingyuanc/GovTech-Technical/models"
)

// Controller for the registration endpoint
func (conn *Connection) HandleSuspend(c *gin.Context) {
	registerBody := models.SuspendBody{}

	// Checks if json body is in correct format
	if err := c.ShouldBind(&registerBody); err != nil {
		// change here if want change error
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

}
