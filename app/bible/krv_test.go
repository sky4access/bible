package bible

import (
	"testing"
)

func TestToSearchableFormat(t *testing.T) {

	var tests = []struct {
		input    []string
		expected []string
	}{
		{
			[]string{"john 1:1, 12, 13", "mark 2:1,5"},
			[]string{"john+1:1", "john+1:12", "john+1:13", "mark+2:1", "mark+2:5"},
		},
		{
			[]string{"john 1:1", "mark 2:1,5,8"},
			[]string{"john+1:1", "mark+2:1", "mark+2:5", "mark+2:8"},
		},
		{
			[]string{"eph 2:13-14, 18, 20"},
			[]string{"eph+2:13-14", "eph+2:18", "eph+2:20"},
		},
	}

	for _, tt := range tests {
		r, err := toSearchableFormat(tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if !stringSlicesEqual(r, tt.expected) {
			t.Errorf("expected %v, got %v", tt.expected, r)
		}
	}

}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
