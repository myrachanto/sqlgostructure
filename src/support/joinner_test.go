package support

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJoins(t *testing.T) {
	testcases := []struct {
		name string
		str  string
		res  string
	}{
		{"ok", "my name is Tony", "my-name-is-Tony"},
		{"empty", "", ""},
	}
	for _, test := range testcases {

		t.Run(test.name, func(t *testing.T) {
			response := Joins(test.str)
			require.EqualValues(t, test.res, response)
		})
	}
}
