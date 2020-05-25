package utils

import (
	"math/rand"
	"regexp"
	"time"
)

const (
	// Alphanumeric charset
	alphaNumCharset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// RFC 3986 Section 2.3 URI Unreserved Characters
	uriUnreservedChars = `^([A-Za-z0-9_.~-])+$`

	// https://gist.github.com/dperini/729294
	urlRegex = `^` +
		// protocol identifier (optional)
		// short syntax // still required
		`(?:(?:(?:https?|ftp):)?\/\/)` +
		// user:pass BasicAuth (optional)
		`(?:\S+(?::\S*)?@)?` +
		`(?:` +
		// IP address dotted notation octets
		// excludes loopback network 0.0.0.0
		// excludes reserved space >= 224.0.0.0
		// excludes network & broadcast addresses
		// (first & last IP address of each class)
		`(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])` +
		`(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}` +
		`(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))` +
		`|` +
		// host & domain names, may end with dot
		// can be replaced by a shortest alternative
		// (?![-_])(?:[-\w\u00a1-\uffff]{0,63}[^-_]\.)+
		`(?:` +
		`(?:` +
		`[a-z0-9\\u00a1-\\uffff]` +
		`[a-z0-9\\u00a1-\\uffff_-]{0,62}` +
		`)?` +
		`[a-z0-9\\u00a1-\\uffff]\.` +
		`)+` +
		// TLD identifier name, may end with dot
		`(?:[a-z\\u00a1-\\uffff]{2,}\.?)` +
		`)` +
		// port number (optional)
		`(?::\d{2,5})?` +
		// resource path (optional)
		`(?:[/?#]\S*)?` +
		`$`
)

// IsValidURL checks if string is a valid url
func IsValidURL(str string) bool {
	valid, err := regexp.MatchString(urlRegex, str)
	return valid && err == nil
}

// IsValidAlias tests whether a string is in the alphanumeric charset
func IsValidAlias(str string) bool {
	valid, err := regexp.MatchString(uriUnreservedChars, str)
	return valid && err == nil
}

// RandString generates a random alphanumeric string of given length
func RandString(length int) string {
	return randStringWithCharset(length, alphaNumCharset)
}

func randStringWithCharset(length int, charset string) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
