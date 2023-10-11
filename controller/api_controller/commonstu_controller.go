package apicontroller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingyuanc/GovTech-Technical/utils"
)

// Controller for the common student endpoint
func (conn *Connection) HandleCommonStu(c *gin.Context) {
	// Get validated data from query
	data, exists := c.Get("teachersParam")
	if !exists {
		// Handle the case where the parameter is not found
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Required teacher query parameter not provided",
		})
		return
	}

	// Access validated data
	teachers, ok := data.([]string)
	if !ok {
		log.Panicf("Error: commonStuController: Unable to cast any to string array, server error")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	students, err := utils.GetUniqueStudentFromTeachersEmail(conn.db, teachers)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"students": students,
	})
}
