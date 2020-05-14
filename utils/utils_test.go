package utils

import (
	"testing"
)

func TestIsAlphaNum(t *testing.T) {
	tests := map[string]struct {
		name string
		str  string
		val  bool
	}{
		"basic":         {str: "hello", val: true},
		"dash":          {str: "-", val: false},
		"period":        {str: ".", val: false},
		"question mark": {str: "?", val: false},
		"backslash":     {str: "?", val: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if IsAlphaNum(tc.str) != tc.val {
				t.Errorf("For '%s', expected %t. Got %t instead.", tc.str, tc.val, IsAlphaNum(tc.str))
			}
		})
	}
}

// For the code coverage why not
func TestRandString(t *testing.T) {
	str := RandString(6)
	if len(str) == 6 && IsAlphaNum(str) {
		return
	}
	t.Errorf("Seriously? How?")
}
