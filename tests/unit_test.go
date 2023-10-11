package tests

import (
	"testing"

	"github.com/mingyuanc/GovTech-Technical/utils"
	"github.com/stretchr/testify/assert"
)

// Tests basic functionality of ExtractEmailFromString function

// Test basic functionality
// Expect extracted email
func TestSplitEmail_Pass(t *testing.T) {
	expected := []string{"studentagnes@gmail.com", "studentmiche@gmail.com"}
	actual := utils.ExtractEmailFromString("Hello students! @studentagnes@gmail.com @studentmiche@gmail.com", "@")
	assert.ElementsMatch(t, expected, actual)
}

// Test basic functionality with no flags and email
// Expect empty array
func TestNoEmail_Pass(t *testing.T) {
	expected := []string{}
	actual := utils.ExtractEmailFromString("Hello Hey everybody!", "@")
	assert.ElementsMatch(t, expected, actual)
}

// Test string with flags but no email
// Expect empty array
func TestHasFlagButInvalidEmail_Pass(t *testing.T) {
	expected := []string{}
	actual := utils.ExtractEmailFromString("Hello Hey @everybody!", "@")
	assert.ElementsMatch(t, expected, actual)
}

// Test basic functionality with no flags but with email
// Expect only emails with flag to be stored
func TestEmailWithoutFlag_Pass(t *testing.T) {
	expected := []string{"testemail@hotmail.com"}
	actual := utils.ExtractEmailFromString("Hello @testemail@hotmail.com Hey studentmiche@gmail.com!", "@")
	assert.ElementsMatch(t, expected, actual)
}
