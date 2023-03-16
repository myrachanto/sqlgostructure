package support

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasher(t *testing.T) {
	testcases := []struct {
		name string
		str  string
		res  string
	}{
		{"ok", "code", "G8wDMmmT4D1tC6R2RGoXeHREZjJ"},
		{"empty", "", ""},
	}
	for _, test := range testcases {

		t.Run(test.name, func(t *testing.T) {
			response := Hasher(test.str)
			require.EqualValues(t, test.res, response)
		})
	}
}
func TestHasher2(t *testing.T) {
	testcases := []struct {
		name string
		str  string
		res  string
	}{
		{"ok", "code", "6pyg7gr1Mhg5kyMrgc5UWb6uGSMUwdHJHYQeb1DWGTbg"},
		{"empty", "", ""},
	}
	for _, test := range testcases {

		t.Run(test.name, func(t *testing.T) {
			response := Hasher2(test.str)
			require.EqualValues(t, test.res, response)
		})
	}
}
