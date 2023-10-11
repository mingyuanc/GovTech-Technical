package utils

import (
	"regexp"
)

// Validates that it is a valid email being used
func IsValidEmail(email string) bool {
	// Regexp pattern for email validation
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, err := regexp.MatchString(emailPattern, email)
	if err != nil {
		return false
	}
	return match
}
