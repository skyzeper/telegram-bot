package utils

import (
	"regexp"
)

// IsValidPhone validates a phone number
func IsValidPhone(phone string) bool {
	// Basic validation: allows formats like +1234567890, 123-456-7890, (123) 456-7890
	re := regexp.MustCompile(`^(\+\d{1,3}[- ]?)?\d{10}$|^(\(\d{3}\)\s?|\d{3}[- ]?)\d{3}[- ]?\d{4}$`)
	return re.MatchString(phone)
}