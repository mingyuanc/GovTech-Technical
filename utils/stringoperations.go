package utils

import (
	"regexp"
	"strings"
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

func ExtractEmailFromString(from string, startsWith string) ([]string, string) {
	var ret = make([]string, 0)
	for _, c := range strings.Split(from, " ") {
		if !strings.HasPrefix(c, startsWith) {
			continue
		}
		email := strings.TrimPrefix(c, startsWith)
		if IsValidEmail(email) {
			ret = append(ret, strings.TrimPrefix(c, startsWith))
			continue
		}
		return ret, email
	}
	return ret, ""
}

func IsInArray(is string, in []string) bool {
	for _, str := range in {
		if is == str {
			return true
		}
	}
	return false
}
