package apicontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingyuanc/GovTech-Technical/models"
	"gorm.io/gorm"
)

type Connection struct {
	db *gorm.DB
}

func NewConnection(db *gorm.DB) *Connection {
	return &Connection{db: db}
}

func (conn *Connection) HandleRegister(c *gin.Context) {
	registerBody := models.Register{}
	if err := c.ShouldBind(&registerBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var teacher models.Teacher
	err := conn.db.Where("email = ?", *registerBody.Teacher).FirstOrCreate(&teacher, models.Teacher{Email: *registerBody.Teacher}).Error
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stuArr := make([]*models.Student, len(registerBody.Students))

	for i, stuEmail := range registerBody.Students {
		var student models.Student
		err := conn.db.Where("email = ?", *stuEmail).FirstOrCreate(&student, models.Student{Email: *stuEmail}).Error
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		stuArr[i] = &student
	}

	conn.db.Model(&teacher).Association("Students").Append(stuArr)

	c.Writer.WriteHeader(http.StatusNoContent)
}
