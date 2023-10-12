package apicontroller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingyuanc/GovTech-Technical/src/utils"
)

// Controller for the common student endpoint
func (conn *Connection) HandleCommonStu(c *gin.Context) {
	conn.logger.Println("api: new common student request")

	// Get validated data from query
	data, exists := c.Get("teachersParam")
	// Another safety check
	if !exists {
		// Handle the case where the parameter is not found
		conn.logger.Println("api/commonstudents: Teacher query parameter not provided")
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Required teacher query parameter not provided",
		})
		return
	}

	// Access validated data
	teachers, ok := data.([]string)
	if !ok {
		conn.logger.Println("api/commonstudents: Unable to cast any to string array, server error")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// Ensure teachers are all present
	for i, teacher := range teachers {
		if !utils.IsTeacherPresent(conn.db, teacher) {
			errStr := fmt.Sprintf("Teacher parameter at index %d is not found: %s", i, teacher)
			conn.logger.Println("api/commonstudents: " + errStr)
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"error": errStr,
			})
			return
		}
	}

	// Runs the query
	students, err := utils.GetCommonStudentFromTeachersEmail(conn.db, teachers)
	if err != nil {
		conn.logger.Println("api/commonstudents: " + err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	// Store the result
	// make() is used so the json returned will show [] instead of null if empty
	var stuEmail = make([]string, 0)
	for _, stu := range students {
		stuEmail = append(stuEmail, stu.Email)
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"students": stuEmail,
	})
}
