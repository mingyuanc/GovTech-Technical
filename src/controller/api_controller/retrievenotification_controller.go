package apicontroller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingyuanc/GovTech-Technical/src/utils"
)

// Controller for the retrieveNotification endpoint
func (conn *Connection) HandleRetrieveNotification(c *gin.Context) {
	conn.logger.Println("api: new retrievefornotifications request")

	// Get validated data from query
	stuArrTmp, notifyExists := c.Get("notify")
	emailTmp, teacherEmailExist := c.Get("teacherEmail")

	// Another safety check
	if !notifyExists || !teacherEmailExist {
		// Handle the case where the parameter is not found
		conn.logger.Println("api/retrievefornotifications: parameter not present after middleware")
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Required parameter not provided",
		})
		return
	}

	// Cast from any to required type
	stuArr, ok := stuArrTmp.([]string)
	if !ok {
		conn.logger.Println("api/retrievefornotifications: Unable to cast any to string array, server error")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	teacherEmail, ok := emailTmp.(string)
	if !ok {
		conn.logger.Println("api/retrievefornotifications: Unable to cast any to string, server error")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	var mentioned = make([]string, 0)
	// check students in mentioned is suspended
	for i, student := range stuArr {

		if !utils.IsStudentPresent(conn.db, student) {
			errStr := fmt.Sprintf("Student parameter at index %d is not found: %s", i, student)
			conn.logger.Println("api/retrievefornotifications: " + errStr)
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"error": errStr,
			})
			return
		}
		if utils.IsStudentSuspended(conn.db, student) {
			continue
		}
		mentioned = append(mentioned, student)
	}

	// Get all student of a teacher
	students, err := utils.GetStudentFromTeacher(conn.db, teacherEmail)
	if err != nil {
		conn.logger.Println("api/retrievefornotifications: " + err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// Store the result
	// make() is used so the json returned will show [] instead of null if empty
	var ret = mentioned
	for _, student := range students {
		if utils.IsInArray(student.Email, mentioned) {
			continue
		}
		if *student.IsSuspended {
			continue
		}
		ret = append(ret, student.Email)
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"recipients": ret,
	})
}
