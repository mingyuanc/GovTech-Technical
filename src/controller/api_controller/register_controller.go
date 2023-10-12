package apicontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingyuanc/GovTech-Technical/src/models"
)

// Controller for the registration endpoint
func (conn *Connection) HandleRegister(c *gin.Context) {
	conn.logger.Println("api: new register request")
	registerBody := models.RegisterBody{}

	// Checks if json body is in correct format
	if err := c.ShouldBind(&registerBody); err != nil {
		// change here if want change error
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//  Check if teacher is present else create new entry
	var teacher models.Teacher
	err := conn.db.Where("email = ?", registerBody.Teacher).FirstOrCreate(&teacher, models.Teacher{Email: registerBody.Teacher}).Error
	if err != nil {
		conn.logger.Println("api/register: " + err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if len(registerBody.Students) == 0 {
		conn.logger.Println("api/register: No student provided")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Must provide at least one student"})
		return
	}

	//  For each student
	//  Check if student is present else create new entry
	//  Then add to stuArr to model the relationship
	stuArr := make([]*models.Student, len(registerBody.Students))
	for i, stuEmail := range registerBody.Students {
		var student models.Student
		err := conn.db.Where("email = ?", stuEmail).FirstOrCreate(&student, models.Student{Email: stuEmail}).Error
		if err != nil {
			conn.logger.Println("api/register: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		stuArr[i] = &student
	}

	conn.db.Model(&teacher).Association("Students").Append(stuArr)

	c.Writer.WriteHeader(http.StatusNoContent)
}
