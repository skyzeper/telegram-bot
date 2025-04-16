package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// FormatPhone formats a phone number to +7(XXX)-XXX-XX-XX
func FormatPhone(phone string) (string, error) {
	// Remove all non-digits
	re := regexp.MustCompile(`\D`)
	digits := re.ReplaceAllString(phone, "")

	// Normalize to 10 digits
	if strings.HasPrefix(digits, "8") {
		digits = "7" + digits[1:]
	} else if strings.HasPrefix(digits, "+7") {
		digits = "7" + digits[2:]
	}

	if len(digits) != 11 {
		return "", fmt.Errorf("invalid phone number length")
	}

	// Format to +7(XXX)-XXX-XX-XX
	return fmt.Sprintf("+7(%s)-%s-%s-%s", digits[1:4], digits[4:7], digits[7:9], digits[9:11]), nil
}

// FormatDate formats a date to "2 January 2006"
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return "В ближайшее время"
	}
	return t.Format("2 January 2006")
}

// FormatTime formats a time to "15:04" or empty for zero
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return "В ближайшее время"
	}
	return t.Format("15:04")
}