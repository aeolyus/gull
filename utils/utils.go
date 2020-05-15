package utils

import (
	"math/rand"
	"regexp"
	"time"
)

const (
	// Alphanumeric charset
	CHARSET = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// RFC 3986 Section 2.3 URI Unreserved Characters
	URIUnreservedChars = `^([A-Za-z0-9_.~-])+$`
)

// Tests whether a string is in the alphanumeric charset
func IsValidAlias(str string) bool {
	valid, err := regexp.MatchString(URIUnreservedChars, str)
	return valid && err == nil
}

// Generate a random alphanumeric string of given length
func RandString(length int) string {
	return randStringWithCharset(length, CHARSET)
}

func randStringWithCharset(length int, charset string) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
