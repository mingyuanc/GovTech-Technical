package api_test

import (
	"bytes"
	"encoding/json"
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
	DB = utils.ConnectTest()
	Router = routes.RunTest(DB)

	m.Run()
}

// Test ping route
func TestPingRoute_200(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/ping", nil)
	Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

type registerBody struct {
	Teacher  string   `json:"teacher"`
	Students []string `json:"students"`
}

// Test register route, correct email
// Expected 204 status code
func TestEmail_204(t *testing.T) {
	email := "teacher@gmail.com"
	students := []string{"stu3@hotmail.com",
		"stu2@u.nus.edu", "stu1@gmail.com"}
	reqBody := &registerBody{Teacher: email, Students: students}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

// Test register route, bad teacher email
// Expected 400 status code
func TestBadTeacherEmail_400(t *testing.T) {
	email := "teacher.com"
	students := []string{"stu3@hotmail.com",
		"stu2@u.nus.edu", "stu1@gmail.com"}
	reqBody := &registerBody{Teacher: email, Students: students}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error": "Key: 'Register.Teacher' Error:Field validation for 'Teacher' failed on the 'email' tag"}`,
		w.Body.String())
}

// Test register route, bad student email for index 0
// Expected 400 status code
func TestOneBadStudentrEmail_400(t *testing.T) {
	email := "teacher@gmail.com"
	students := []string{"stu3.com",
		"stu2@u.nus.edu", "stu1@gmail.com"}
	reqBody := &registerBody{Teacher: email, Students: students}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{
		"error": "Key: 'Register.Students[0]' Error:Field validation for 'Students[0]' failed on the 'email' tag"
	}`,
		w.Body.String())
}

// Test register route, empty student array
// Expected 400 status code
func TestNoStudent_400(t *testing.T) {
	email := "teacher@gmail.com"
	students := []string{}
	reqBody := &registerBody{Teacher: email, Students: students}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{
		"error": "Must provide at least one student"
	}`,
		w.Body.String())
}

// Test register route, teacher to student association
// Expected to exit without error
func TestTeacherToMutiStudentAssociation_Pass(t *testing.T) {
	email := "associator@gmail.com"
	students := []string{"associated1@gmail.com"}
	reqBody := &registerBody{Teacher: email, Students: students}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)

	var teacher models.Teacher
	DB.Preload("Students").Where("email = ?", email).First(&teacher)
	assert.Equal(t, teacher.Students[0].Email, "associated1@gmail.com")

	// add another student to teacher
	students = []string{"associated2@gmail.com"}
	reqBody = &registerBody{Teacher: email, Students: students}
	json.NewEncoder(payloadBuf).Encode(reqBody)

	req, _ = http.NewRequest("POST", "/api/register", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)

	DB.Preload("Students").Where("email = ?", email).First(&teacher)
	assert.Equal(t, teacher.Students[0].Email, "associated1@gmail.com")
	assert.Equal(t, teacher.Students[1].Email, "associated2@gmail.com")
}

// Test register route, teacher to student association
// Expected to exit without error
func TestStudentToMutiTeacherAssociation_Pass(t *testing.T) {
	email := "associated1@gmail.com"
	students := []string{"associator@gmail.com"}
	reqBody := &registerBody{Teacher: email, Students: students}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)

	var student models.Student
	DB.Preload("Teachers").Where("email = ?", "associator@gmail.com").First(&student)
	assert.Equal(t, student.Teachers[0].Email, "associated1@gmail.com")

	// add another teacher to student
	email = "associated2@gmail.com"
	reqBody = &registerBody{Teacher: email, Students: students}
	json.NewEncoder(payloadBuf).Encode(reqBody)

	req, _ = http.NewRequest("POST", "/api/register", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)

	DB.Preload("Teachers").Where("email = ?", "associator@gmail.com").First(&student)
	assert.Equal(t, student.Teachers[0].Email, "associated1@gmail.com")
	assert.Equal(t, student.Teachers[1].Email, "associated2@gmail.com")
}
