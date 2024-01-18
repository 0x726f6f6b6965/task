package helper

import "strings"

// IsEmpty - check if string is empty
func IsEmpty(s string) bool {
	return strings.Trim(s, " ") == ""
}
