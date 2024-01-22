package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMD5(t *testing.T) {

	type TestCase struct {
		In  string
		Out string
	}
	testCases := []TestCase{
		{
			In:  "2/00/token",
			Out: "3d9070e0",
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.Out, MD5(tc.In))
	}
}
