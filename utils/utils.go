package utils

import (
	"math/rand"
	"strings"
	"time"
)

// Alphanumeric charset
const CHARSET = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Tests whether a string is in the alphanumeric charset
func IsAlphaNum(str string) bool {
	for _, r := range []rune(str) {
		if !strings.ContainsRune(CHARSET, r) {
			return false
		}
	}
	return true
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
