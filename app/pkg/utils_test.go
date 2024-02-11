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

func TestF64ToS(t *testing.T) {
	type TestCase struct {
		In  float64
		Out string
	}
	testCases := []TestCase{
		{
			In:  2.0,
			Out: "2",
		},
		{
			In:  2.1,
			Out: "2",
		},
		{
			In:  2.5,
			Out: "2",
		},
		{
			In:  2.9,
			Out: "3",
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.Out, F64ToS(&tc.In))
	}
}

func TestF64To1DecimalF64(t *testing.T) {
	type TestCase struct {
		In  float64
		Out float64
	}
	testCases := []TestCase{
		{
			In:  2.0,
			Out: 2.0,
		},
		{
			In:  2.1,
			Out: 2.1,
		},
		{
			In:  2.5,
			Out: 2.5,
		},
		{
			In:  2.9,
			Out: 2.9,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.Out, F64To1DecimalF64(&tc.In))
	}
}

func TestStringToInt(t *testing.T) {

	type TestCase struct {
		In  string
		Out int64
	}
	testCases := []TestCase{
		{
			In:  "2",
			Out: 2,
		},
		{
			In:  "2.1",
			Out: 2,
		},
		{
			In:  "2.5",
			Out: 2,
		},
		{
			In:  "2.9",
			Out: 2,
		},
		{
			In:  "2.9.9",
			Out: 2,
		},
		{
			In:  "",
			Out: 0,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.Out, StringToInt(tc.In))
	}
}

func TestSToF32(t *testing.T) {

	type TestCase struct {
		In  string
		Out float32
	}
	testCases := []TestCase{
		{
			In:  "2",
			Out: 2,
		},
		{
			In:  "2.1",
			Out: 2.1,
		},
		{
			In:  "2.5",
			Out: 2.5,
		},
		{
			In:  "2.9",
			Out: 2.9,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.Out, SToF32(tc.In))
	}
}

func TestSToF64(t *testing.T) {

	type TestCase struct {
		In  string
		Out float64
	}
	testCases := []TestCase{
		{
			In:  "2",
			Out: 2,
		},
		{
			In:  "2.1",
			Out: 2.1,
		},
		{
			In:  "2.5",
			Out: 2.5,
		},
		{
			In:  "2.9",
			Out: 2.9,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.Out, SToF64(tc.In))
	}
}

func TestTakeFirst(t *testing.T) {

	type TestCase struct {
		In  string
		N   int
		Out string
	}
	testCases := []TestCase{
		{
			In:  "2",
			N:   2,
			Out: "2",
		},
		{
			In:  "2.1",
			N:   2,
			Out: "2.",
		},
		{
			In:  "2.5",
			N:   2,
			Out: "2.",
		},
		{
			In:  "2.9",
			N:   2,
			Out: "2.",
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.Out, TakeFirst(tc.In, tc.N))
	}
}
