package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/mingyuanc/GovTech-Technical/models"
	"github.com/stretchr/testify/assert"
)

func CommonStuSetUp() {

	students := []*models.Student{
		{Email: "testcommon1@gmail.com"},
		{Email: "testcommon2@gmail.com"},
	}

	unique := models.Student{Email: "testunique@gmail.com"}
	ken := models.Teacher{
		Email: "testken@gmail.com",
	}
	ming := models.Teacher{
		Email: "testming@gmail.com",
	}
	noCommon := models.Teacher{
		Email: "testnocommon@gmail.com",
	}

	DB.Create(&ken)
	DB.Create(&ming)
	DB.Create(&noCommon)
	for _, student := range students {
		DB.Create(student)
		DB.Model(student).Association("Teachers").Append([]models.Teacher{ken, ming})
	}
	DB.Model(&ken).Association("Students").Append(&unique)
}

// Test get common student, expected 3 students as per given test case
// Expected 200 status code
func TestProvidedCase1_200(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/commonstudents?teacher=testken%40gmail.com", nil)
	Router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := map[string]interface{}{
		"students": []interface{}{"testcommon1@gmail.com", "testcommon2@gmail.com", "testunique@gmail.com"},
	}
	var data map[string]interface{}

	// Use the JSON unmarshal function to decode the binary data into a map
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Fail()
	}

	assert.True(t, reflect.DeepEqual(expected, data))
}

// Test get common student, expected 2 students as per given test case
// Expected 200 status code
func TestProvidedCase2_200(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/commonstudents?teacher=testken%40gmail.com&teacher=testming%40gmail.com", nil)
	Router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := map[string]interface{}{
		"students": []interface{}{"testcommon1@gmail.com", "testcommon2@gmail.com"},
	}
	var data map[string]interface{}

	// Use the JSON unmarshal function to decode the binary data into a map
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Fail()
	}

	assert.True(t, reflect.DeepEqual(expected, data))
}

// Test get common student, expected no common student
// Expected 200 status code
func TestNoCommon_200(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/commonstudents?teacher=testken%40gmail.com&teacher=testnocommon%40gmail.com", nil)
	Router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := map[string]interface{}{
		"students": []interface{}{},
	}
	var data map[string]interface{}

	// Use the JSON unmarshal function to decode the binary data into a map
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Fail()
	}

	assert.True(t, reflect.DeepEqual(expected, data))
}

// Test no parameters
// both expected to have 400 response
func TestNoParam_400(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/commonstudents", nil)
	Router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error": "Required teacher query parameter not provided"}`, w.Body.String())

}

// Test get teacher invalid email
// both expected to have 400 response
func TestTeacherBadEmail_400(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/commonstudents?teacher=testNotFound%40gmail", nil)
	Router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error": "Teacher parameter at index 0 is an invalid email: testNotFound@gmail"}`, w.Body.String())
}

// Test get no such teacher found
// both expected to have 400 response
func TestTeacherNotFound_400(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/commonstudents?teacher=testNotFound%40gmail.com", nil)
	Router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error": "Teacher parameter at index 0 is not found: testNotFound@gmail.com"}`, w.Body.String())
}
