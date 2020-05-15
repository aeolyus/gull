package utils

import (
	"testing"
)

func TestIsValidAlias(t *testing.T) {
	tests := map[string]struct {
		name string
		str  string
		val  bool
	}{
		"LOWERCASE":     {str: "lowercase", val: true},
		"UPPERCASE":     {str: "UPPERCASE", val: true},
		"underscore":    {str: "_", val: true},
		"period":        {str: ".", val: true},
		"dash":          {str: "-", val: true},
		"tilde":         {str: "~", val: true},
		"adsf":          {str: "a_S-d.f~", val: true},
		"question mark": {str: "?", val: false},
		"backslash":     {str: "\\", val: false},
		"asterix":       {str: "*", val: false},
		"ampersand":     {str: "&", val: false},
		"space":         {str: " ", val: false},
		"percent":       {str: "%", val: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res := IsValidAlias(tc.str)
			if res != tc.val {
				t.Errorf("For '%s', expected %t. Got %t instead.", tc.str, tc.val, res)
			}
		})
	}
}

func TestIsValidURL(t *testing.T) {
	tests := map[string]struct {
		name string
		url  string
		val  bool
	}{
		// Should match
		"http://foo.com/blah_blah":                          {url: "http://foo.com/blah_blah", val: true},
		"http://foo.com/blah_blah/":                         {url: "http://foo.com/blah_blah/", val: true},
		"http://foo.com/blah_blah_(wikipedia)":              {url: "http://foo.com/blah_blah_(wikipedia)", val: true},
		"http://foo.com/blah_blah_(wikipedia)_(again)":      {url: "http://foo.com/blah_blah_(wikipedia)_(again)", val: true},
		"http://www.example.com/wpstyle/?p=364":             {url: "http://www.example.com/wpstyle/?p=364", val: true},
		"https://www.example.com/foo/?bar=baz&inga=42&quux": {url: "https://www.example.com/foo/?bar=baz&inga=42&quux", val: true},
		"http://userid:password@example.com:8080":           {url: "http://userid:password@example.com:8080", val: true},
		"http://userid:password@example.com:8080/":          {url: "http://userid:password@example.com:8080/", val: true},
		"http://userid@example.com":                         {url: "http://userid@example.com", val: true},
		"http://userid@example.com/":                        {url: "http://userid@example.com/", val: true},
		"http://userid@example.com:8080":                    {url: "http://userid@example.com:8080", val: true},
		"http://userid@example.com:8080/":                   {url: "http://userid@example.com:8080/", val: true},
		"http://userid:password@example.com":                {url: "http://userid:password@example.com", val: true},
		"http://userid:password@example.com/":               {url: "http://userid:password@example.com/", val: true},
		"http://142.42.1.1/":                                {url: "http://142.42.1.1/", val: true},
		"http://142.42.1.1:8080/":                           {url: "http://142.42.1.1:8080/", val: true},
		"http://foo.com/blah_(wikipedia)#cite-1":            {url: "http://foo.com/blah_(wikipedia)#cite-1", val: true},
		"http://foo.com/blah_(wikipedia)_blah#cite-1":       {url: "http://foo.com/blah_(wikipedia)_blah#cite-1", val: true},
		"http://foo.com/unicode_(✪)_in_parens":              {url: "http://foo.com/unicode_(✪)_in_parens", val: true},
		"http://foo.com/(something)?after=parens":           {url: "http://foo.com/(something)?after=parens", val: true},
		"http://code.google.com/events/#&product=browser":   {url: "http://code.google.com/events/#&product=browser", val: true},
		"http://j.mp":       {url: "http://j.mp", val: true},
		"ftp://foo.bar/baz": {url: "ftp://foo.bar/baz", val: true},
		"http://foo.bar/?q=Test%20URL-encoded%20stuff": {url: "http://foo.bar/?q=Test%20URL-encoded%20stuff", val: true},
		"http://1337.net":              {url: "http://1337.net", val: true},
		"http://a.b-c.de":              {url: "http://a.b-c.de", val: true},
		"http://223.255.255.254":       {url: "http://223.255.255.254", val: true},
		"https://foo_bar.example.com/": {url: "https://foo_bar.example.com/", val: true},
		"http://www.foo.bar./":         {url: "http://www.foo.bar.", val: true},
		"http://a.b--c.de/":            {url: "http://a.b--c.de/", val: true},
		// Should not match
		"http://":    {url: "http://", val: false},
		"http://.":   {url: "http://.", val: false},
		"http://..":  {url: "http://..", val: false},
		"http://../": {url: "http://../", val: false},
		"http://?":   {url: "http://?", val: false},
		"http://??":  {url: "http://??", val: false},
		"http://??/": {url: "http://??/", val: false},
		"http://#":   {url: "http://#", val: false},
		"http://##":  {url: "http://##", val: false},
		"http://##/": {url: "http://##/", val: false},
		"http://foo.bar?q=Spaces should be encoded": {url: "http://foo.bar?q=Spaces should be encoded", val: false},
		"//":                              {url: "//", val: false},
		"//a":                             {url: "//a", val: false},
		"///a":                            {url: "///a", val: false},
		"///":                             {url: "///", val: false},
		"http:///a":                       {url: "http:///a", val: false},
		"foo.com":                         {url: "foo.com", val: false},
		"rdar://1234":                     {url: "rdar://1234", val: false},
		"h://test":                        {url: "h://test", val: false},
		"http:// shouldfail.com":          {url: "http:// shouldfail.com", val: false},
		":// should fail":                 {url: ":// should fail", val: false},
		"http://foo.bar/foo(bar)baz quux": {url: "http://foo.bar/foo(bar)baz quux", val: false},
		"ftps://foo.bar/":                 {url: "ftps://foo.bar/", val: false},
		"http://-error-.invalid/":         {url: "http://-error-.invalid/", val: false},
		"http://-a.b.co":                  {url: "http://-a.b.co", val: false},
		"http://a.b-.co":                  {url: "http://a.b-.co", val: false},
		"http://0.0.0.0":                  {url: "http://0.0.0.0", val: false},
		"http://10.1.1.0":                 {url: "http://10.1.1.0", val: false},
		"http://224.1.1.1":                {url: "http://224.1.1.1", val: false},
		"http://1.1.1.1.1":                {url: "http://1.1.1.1.1", val: false},
		"http://3628126748":               {url: "http://3628126748", val: false},
		"http://.www.foo.bar/":            {url: "http://.www.foo.bar/", val: false},
		"http://.www.foo.bar./":           {url: "http://.www.foo.bar./", val: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res := IsValidURL(tc.url)
			if res != tc.val {
				t.Errorf("For '%s', expected %t. Got %t instead.", tc.url, tc.val, res)
			}
		})
	}
}

// For the code coverage why not
func TestRandString(t *testing.T) {
	str := RandString(6)
	if len(str) == 6 && IsValidAlias(str) {
		return
	}
	t.Errorf("Seriously? How?")
}
