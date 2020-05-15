package utils

import (
	"testing"
)

func TestValidAlias(t *testing.T) {
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
			if IsValidAlias(tc.str) != tc.val {
				t.Errorf("For '%s', expected %t. Got %t instead.", tc.str, tc.val, IsValidAlias(tc.str))
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
