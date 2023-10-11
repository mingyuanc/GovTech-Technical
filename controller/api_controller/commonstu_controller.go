package apicontroller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingyuanc/GovTech-Technical/utils"
)

// Controller for the common student endpoint
func (conn *Connection) HandleCommonStu(c *gin.Context) {
	// Get validated data from query
	data, exists := c.Get("teachersParam")
	// Another safety check
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

	// Ensure teachers are all present
	for i, teacher := range teachers {
		if !utils.IsTeacherPresent(conn.db, teacher) {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Teacher parameter at index %d is not found: %s", i, teacher),
			})
			return
		}
	}

	// Runs the query
	students, err := utils.GetCommonStudentFromTeachersEmail(conn.db, teachers)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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
