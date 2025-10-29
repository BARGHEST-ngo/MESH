package hashes

import (
	"regexp"
)

// ValidateMD5 returns true if the specified string is a valid MD5 hash.
func ValidateMD5(target string) bool {
	rxp := regexp.MustCompile(`[a-fA-F0-9]{32}`)
	return rxp.MatchString(target)
}

// ValidateSHA1 returns true if the specified string is a valid SHA1 hash.
func ValidateSHA1(target string) bool {
	rxp := regexp.MustCompile(`[a-fA-F0-9]{40}`)
	return rxp.MatchString(target)
}

// ValidateSHA256 returns true if the specified string is a valid SHA256 hash.
func ValidateSHA256(target string) bool {
	rxp := regexp.MustCompile(`[a-fA-F0-9]{64}`)
	return rxp.MatchString(target)
}

// ValidateSHA512 returns true if the specified string is a valid SHA512 hash.
func ValidateSHA512(target string) bool {
	rxp := regexp.MustCompile(`[a-fA-F0-9]{128}`)
	return rxp.MatchString(target)
}
