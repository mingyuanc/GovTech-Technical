package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mingyuanc/GovTech-Technical/models"
	"github.com/stretchr/testify/assert"
)

type registerBody struct {
	Teacher  string   `json:"teacher"`
	Students []string `json:"students"`
}

// Test register route, correct email
// Expected 204 status code
func TestEmail_204(t *testing.T) {
	email := "testTeacher@gmail.com"
	students := []string{"teststu3@hotmail.com",
		"teststu2@u.nus.edu", "teststu1@gmail.com"}
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
	email := "testteacher.com"
	students := []string{"teststu3@hotmail.com",
		"teststu2@u.nus.edu", "teststu1@gmail.com"}
	reqBody := &registerBody{Teacher: email, Students: students}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error": "Key: 'RegisterBody.Teacher' Error:Field validation for 'Teacher' failed on the 'email' tag"}`,
		w.Body.String())
}

// Test register route, bad student email for index 0
// Expected 400 status code
func TestOneBadStudentEmail_400(t *testing.T) {
	email := "testTeacher@gmail.com"
	students := []string{"teststu3.com",
		"teststu2@u.nus.edu", "teststu1@gmail.com"}
	reqBody := &registerBody{Teacher: email, Students: students}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{
		"error": "Key: 'RegisterBody.Students[0]' Error:Field validation for 'Students[0]' failed on the 'email' tag"
	}`,
		w.Body.String())
}

// Test register route, empty student array
// Expected 400 status code
func TestNoStudent_400(t *testing.T) {
	email := "testTeacher@gmail.com"
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
func TestTeacherToMultiStudentAssociation_Pass(t *testing.T) {
	email := "testassociator@gmail.com"
	students := []string{"testassociated1@gmail.com"}
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
	assert.Equal(t, "testassociated1@gmail.com", teacher.Students[0].Email)

	// add another student to teacher
	students = []string{"testassociated2@gmail.com"}
	reqBody = &registerBody{Teacher: email, Students: students}
	json.NewEncoder(payloadBuf).Encode(reqBody)

	req, _ = http.NewRequest("POST", "/api/register", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)

	DB.Preload("Students").Where("email = ?", email).First(&teacher)
	assert.Equal(t, "testassociated1@gmail.com", teacher.Students[0].Email)
	assert.Equal(t, "testassociated2@gmail.com", teacher.Students[1].Email)
}

// Test register route, teacher to student association
// Expected to exit without error
func TestStudentToMultiTeacherAssociation_Pass(t *testing.T) {
	email := "testassociated1@gmail.com"
	students := []string{"testassociator@gmail.com"}
	reqBody := &registerBody{Teacher: email, Students: students}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)

	var student models.Student
	DB.Preload("Teachers").Where("email = ?", "testassociator@gmail.com").First(&student)
	assert.Equal(t, "testassociated1@gmail.com", student.Teachers[0].Email)

	// add another teacher to student
	email = "testassociated2@gmail.com"
	reqBody = &registerBody{Teacher: email, Students: students}
	json.NewEncoder(payloadBuf).Encode(reqBody)

	req, _ = http.NewRequest("POST", "/api/register", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)

	DB.Preload("Teachers").Where("email = ?", "testassociator@gmail.com").First(&student)
	assert.Equal(t, "testassociated1@gmail.com", student.Teachers[0].Email)
	assert.Equal(t, "testassociated2@gmail.com", student.Teachers[1].Email)
}
