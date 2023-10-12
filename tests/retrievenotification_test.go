package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/mingyuanc/GovTech-Technical/src/models"
	"github.com/mingyuanc/GovTech-Technical/src/utils"
	"github.com/stretchr/testify/assert"
)

func RetrieveNotificationSetup() {

	students := []*models.Student{
		{Email: "testsuspend1@gmail.com"},
		{Email: "testsuspend2@gmail.com"},
		{Email: "testnotsuspend1@gmail.com"},
		{Email: "testnotsuspend2@gmail.com"},
	}

	notRelatedStudent := models.Student{
		Email: "testnotrelated@hotmail.com",
	}
	DB.Create(&notRelatedStudent)
	suspendTeacher := models.Teacher{
		Email: "testsuspendteacher@gmail.com",
	}

	DB.Create(&suspendTeacher)
	for _, student := range students {
		DB.Create(student)
		DB.Model(student).Association("Teachers").Append(&suspendTeacher)
	}
}

type retrieveNotificationBody struct {
	Teacher      string `json:"teacher"`
	Notification string `json:"notification"`
}

func TestNoSuspended_200(t *testing.T) {
	reqBody := &retrieveNotificationBody{Teacher: "testsuspendteacher@gmail.com",
		Notification: "Hey @testnotrelated@hotmail.com you should appear"}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/retrievefornotifications", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := map[string]interface{}{
		"recipients": []interface{}{"testnotrelated@hotmail.com", "testsuspend1@gmail.com", "testsuspend2@gmail.com", "testnotsuspend1@gmail.com", "testnotsuspend2@gmail.com"},
	}
	var data map[string]interface{}

	// Use the JSON unmarshal function to decode the binary data into a map
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Fail()
	}
	assert.True(t, reflect.DeepEqual(expected, data))
}

func TestEmptyNotification_200(t *testing.T) {
	reqBody := &retrieveNotificationBody{Teacher: "testsuspendteacher@gmail.com"}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/retrievefornotifications", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := map[string]interface{}{
		"recipients": []interface{}{"testsuspend1@gmail.com", "testsuspend2@gmail.com", "testnotsuspend1@gmail.com", "testnotsuspend2@gmail.com"},
	}
	var data map[string]interface{}

	// Use the JSON unmarshal function to decode the binary data into a map
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Fail()
	}
	assert.True(t, reflect.DeepEqual(expected, data))
}

func TestNoDuplicateWithStudentInNotification_200(t *testing.T) {
	reqBody := &retrieveNotificationBody{Teacher: "testsuspendteacher@gmail.com",
		Notification: "Hey @testsuspend1@gmail.com and @testsuspend2@gmail.com you should only appear once"}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/retrievefornotifications", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := map[string]interface{}{
		"recipients": []interface{}{"testsuspend1@gmail.com", "testsuspend2@gmail.com", "testnotsuspend1@gmail.com", "testnotsuspend2@gmail.com"},
	}
	var data map[string]interface{}

	// Use the JSON unmarshal function to decode the binary data into a map
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Fail()
	}
	assert.True(t, reflect.DeepEqual(expected, data))
}

func TestSuspended_200(t *testing.T) {
	// Suspend student
	utils.SuspendStudent(DB, "testsuspend1@gmail.com")
	utils.SuspendStudent(DB, "testsuspend2@gmail.com")

	reqBody := &retrieveNotificationBody{Teacher: "testsuspendteacher@gmail.com",
		Notification: "Hey @testsuspend1@gmail.com you should not appear but @testnotrelated@hotmail.com should appear"}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/retrievefornotifications", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := map[string]interface{}{
		"recipients": []interface{}{"testnotrelated@hotmail.com", "testnotsuspend1@gmail.com", "testnotsuspend2@gmail.com"},
	}
	var data map[string]interface{}

	// Use the JSON unmarshal function to decode the binary data into a map
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Fail()
	}
	assert.True(t, reflect.DeepEqual(expected, data))
}

func TestSuspended2_200(t *testing.T) {
	// Suspend student
	utils.SuspendStudent(DB, "testsuspend1@gmail.com")
	utils.SuspendStudent(DB, "testsuspend2@gmail.com")

	reqBody := &retrieveNotificationBody{Teacher: "testsuspendteacher@gmail.com",
		Notification: "Hey @testsuspend1@gmail.com you should not appear"}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/retrievefornotifications", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := map[string]interface{}{
		"recipients": []interface{}{"testnotsuspend1@gmail.com", "testnotsuspend2@gmail.com"},
	}
	var data map[string]interface{}

	// Use the JSON unmarshal function to decode the binary data into a map
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Fail()
	}
	assert.True(t, reflect.DeepEqual(expected, data))
}

func TestBadMentionedEmail_400(t *testing.T) {

	reqBody := &retrieveNotificationBody{Teacher: "testsuspendteacher@gmail.com",
		Notification: "Hey @testsuspend1.com you should not appear"}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/retrievefornotifications", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	assert.JSONEq(t, `{"error": "Student mentioned has an invalid email: testsuspend1.com"}`, w.Body.String())
}

func TestBadTeacherEmailRetrieval_400(t *testing.T) {

	reqBody := &retrieveNotificationBody{Teacher: "testsuspendteacher",
		Notification: "Hey @testsuspend1.com you should not appear"}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/retrievefornotifications", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	assert.JSONEq(t, `{"error": "Key: 'RetrieveNotificationBody.Teacher' Error:Field validation for 'Teacher' failed on the 'email' tag"}`, w.Body.String())
}
