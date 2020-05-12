package utils

import (
	"math/rand"
	"net/url"
	"time"
)

// Charset for random string generator
const Charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Tests a string to determine if it is a well-structured url or not
func IsValidUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Generate a random alphanumeric string of given length
func RandString(length int) string {
	return randStringWithCharset(length, Charset)
}

func randStringWithCharset(length int, charset string) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
