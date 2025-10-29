package text

import (
	"strings"
)

// Contains returns true if the specified text (inclusive of new-lines)
// contains the specified string (case-sensitive).
func Contains(text, pattern string) bool {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.Contains(line, pattern) {
			return true
		}
	}

	return false
}

// ContainsNoCase returns true if the specified text (inclusive of
// new-lines) contains the specified string (case-sensitive).
func ContainsNoCase(text, pattern string) bool {
	pattern = strings.ToLower(pattern)
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.ToLower(line)
		if strings.Contains(line, pattern) {
			return true
		}
	}

	return false
}
