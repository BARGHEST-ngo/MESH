package slice

import (
	"strings"
)

func stringInSlice(pattern string, slice []string, caseSensitive bool) bool {
	if !caseSensitive {
		pattern = strings.ToLower(pattern)
	}

	for _, entry := range slice {
		if !caseSensitive {
			entry = strings.ToLower(entry)
		}
		if entry == pattern {
			return true
		}
	}

	return false
}

// Contains returns true if the specified string slice contains the
// specified string (case-sensitive).
func Contains(slice []string, pattern string) bool {
	return stringInSlice(pattern, slice, true)
}

// ContainsNoCase returns true if the specified string slice contains the
// specified string (case-insensitive).
func ContainsNoCase(slice []string, pattern string) bool {
	return stringInSlice(pattern, slice, false)
}
