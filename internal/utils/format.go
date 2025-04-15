package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func FormatPhone(phone string) string {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	re := regexp.MustCompile(`^(?:\+7|8|7)?(\d{10})$`)
	matches := re.FindStringSubmatch(phone)
	if len(matches) != 2 {
		return ""
	}
	digits := matches[1]
	return fmt.Sprintf("+7(%s)-%s-%s-%s", digits[:3], digits[3:6], digits[6:8], digits[8:10])
}
