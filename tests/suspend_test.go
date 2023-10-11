package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mingyuanc/GovTech-Technical/models"
	"github.com/stretchr/testify/assert"
)

func SuspendSetUp() {

	students := []*models.Student{
		{Email: "testsuspend1@gmail.com"},
		{Email: "testsuspend2@gmail.com"},
		{Email: "testnotsuspend1@gmail.com"},
		{Email: "testnotsuspend2@gmail.com"},
	}

	suspendTeacher := models.Teacher{
		Email: "testsuspendteacher@gmail.com",
	}

	DB.Create(&suspendTeacher)
	for _, student := range students {
		DB.Create(student)
		DB.Model(student).Association("Teachers").Append(&suspendTeacher)
	}
}

type suspendBody struct {
	Student string `json:"student"`
}

// Test suspend route, correct email
// Expected 204 status code
func TestSuspend_204(t *testing.T) {
	reqBody := &suspendBody{Student: "testsuspend1@gmail.com"}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/suspend", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)
}

// Test suspend route, bad student email
// Expected 400 status code
func TestBadStudentEmail_400(t *testing.T) {
	reqBody := &suspendBody{Student: "testsuspend1.com"}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/suspend", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error": "Key: 'SuspendBody.Student' Error:Field validation for 'Student' failed on the 'email' tag"}`,
		w.Body.String())
}

// Test suspend route, student not found
// Expected 400 status code
func TestStudentNotFound_400(t *testing.T) {
	reqBody := &suspendBody{Student: "testsuspend3@gmail.com"}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/suspend", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error": "Student not found: testsuspend3@gmail.com"}`,
		w.Body.String())
}
