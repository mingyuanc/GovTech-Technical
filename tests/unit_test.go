package tests

import (
	"testing"

	"github.com/mingyuanc/GovTech-Technical/src/utils"
	"github.com/stretchr/testify/assert"
)

// Tests basic functionality of ExtractEmailFromString function

// Test basic functionality
// Expect extracted email
func TestSplitEmail_Pass(t *testing.T) {
	expected := []string{"studentagnes@gmail.com", "studentmiche@gmail.com"}
	actual, notValid := utils.ExtractEmailFromString("Hello students! @studentagnes@gmail.com @studentmiche@gmail.com", "@")
	if notValid != "" {
		t.Fail()
	}
	assert.ElementsMatch(t, expected, actual)
}

// Test basic functionality with no flags and email
// Expect empty array
func TestNoEmail_Pass(t *testing.T) {
	expected := []string{}
	actual, notValid := utils.ExtractEmailFromString("Hello Hey everybody!", "@")
	if notValid != "" {
		t.Fail()
	}
	assert.ElementsMatch(t, expected, actual)
}

// Test string with flags but no email
// Expect to fail
func TestHasFlagButInvalidEmail_NotValid(t *testing.T) {
	_, notValid := utils.ExtractEmailFromString("Hello Hey @everybody!", "@")
	if notValid != "" {
		assert.Equal(t, "everybody!", notValid)
		return
	}
	t.Fail()
}

// Test basic functionality with no flags but with email
// Expect only emails with flag to be stored
func TestEmailWithoutFlag_Pass(t *testing.T) {
	expected := []string{"testemail@hotmail.com"}
	actual, notValid := utils.ExtractEmailFromString("Hello @testemail@hotmail.com Hey studentmiche@gmail.com!", "@")
	if notValid != "" {
		t.Fail()
	}
	assert.ElementsMatch(t, expected, actual)
}
