package internal_test

import (
	"github.com/avnovoselov/live_debugger/internal/server/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalizeLocation(t *testing.T) {
	type TestCase struct {
		name string
		in   string
		out  string
	}

	testCases := []TestCase{
		{
			name: "Basic",
			in:   "foo",
			out:  "/foo",
		},
		{
			name: "Trim leading slash",
			in:   "///foo",
			out:  "/foo",
		},
		{
			name: "Keep not leading slash",
			in:   "foo/bar",
			out:  "/foo/bar",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := internal.NormalizeLocation(testCase.in)
			assert.Equal(t, testCase.out, r)
		})
	}
}
