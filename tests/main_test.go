package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mingyuanc/GovTech-Technical/models"
	"github.com/mingyuanc/GovTech-Technical/routes"
	"github.com/mingyuanc/GovTech-Technical/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var Router *gin.Engine
var DB *gorm.DB

// Creates the required variable and starts the test
func TestMain(m *testing.M) {
	// Sets up global variable and required data in database
	DB = utils.Connect()
	Router = routes.RunTest(DB)
	CommonStuSetUp()
	SuspendSetUp()

	m.Run()

	// clean up
	var delTeachers []models.Teacher
	DB.Where("email LIKE ?", "test%").Find(&delTeachers)

	// Step 2: Delete many-to-many associations
	for _, teacher := range delTeachers {
		DB.Model(&teacher).Association("Students").Unscoped().Clear()
	}

	// Step 3: Delete the teachers
	DB.Unscoped().Delete(&delTeachers)
	DB.Unscoped().Delete(&models.Student{}, "email LIKE ? ", "test%")
}

// Test ping route
func TestPingRoute_200(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/ping", nil)
	Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
