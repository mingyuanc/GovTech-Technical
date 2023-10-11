package apicontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingyuanc/GovTech-Technical/models"
	"github.com/mingyuanc/GovTech-Technical/utils"
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
	studentEmail := registerBody.Student
	if !utils.IsStudentPresent(conn.db, studentEmail) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Student not found: " + studentEmail})
		return
	}
	if utils.SuspendStudent(conn.db, studentEmail) != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to suspend student"})
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
