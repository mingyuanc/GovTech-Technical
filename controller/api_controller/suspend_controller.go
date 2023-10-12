package apicontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingyuanc/GovTech-Technical/models"
	"github.com/mingyuanc/GovTech-Technical/utils"
)

// Controller for the registration endpoint
func (conn *Connection) HandleSuspend(c *gin.Context) {
	conn.logger.Println("api: new suspend request")
	suspendBody := models.SuspendBody{}

	// Checks if json body is in correct format
	if err := c.ShouldBind(&suspendBody); err != nil {
		// change here if want change error
		conn.logger.Println("api/suspend: Unable to bind to body")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	studentEmail := suspendBody.Student
	if !utils.IsStudentPresent(conn.db, studentEmail) {
		conn.logger.Println("api/suspend: student not present")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Student not found: " + studentEmail})
		return
	}
	if utils.SuspendStudent(conn.db, studentEmail) != nil {
		conn.logger.Println("api/suspend: Unable to suspend student")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to suspend student"})
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
