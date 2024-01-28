package pkg

import (
	"hash/fnv"
	"strconv"
	"strings"
)

func F64ToS(f *float64) string {
	return strconv.FormatFloat(*f, 'f', 0, 64)
}

func F64NumberToK(num *float64) string {
	if num == nil {
		return "0"
	}

	if *num < 1000 {
		return strconv.FormatFloat(*num, 'f', -1, 64)
	}

	if *num < 1000000 {
		return strconv.FormatFloat(*num/1000, 'f', 1, 64) + "k"
	}

	return strconv.FormatFloat(*num/1000000, 'f', 1, 64) + "m"
}

// float64 to one decimal place
func F64To1DecimalF64(num *float64) float64 {
	if num == nil {
		return 0
	}
	return float64(int(*num*10)) / 10
}

func StringToInt(s string) int64 {
	parts := strings.Split(s, ".")
	if len(parts) == 0 {
		return 0
	}

	s = parts[0]
	i, _ := strconv.Atoi(s)
	return int64(i)
}

func SToF32(s string) float32 {
	i, _ := strconv.ParseFloat(s, 32)
	return float32(i)
}
func SToF64(s string) float64 {
	i, _ := strconv.ParseFloat(s, 64)
	return float64(i)
}

func TakeFirst(s string, n int) string {
	if len(s) < n {
		return s
	}
	return s[:n]
}

func MD5(s string) string {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return s
	}
	u := h.Sum32()
	return strconv.FormatUint(uint64(u), 16)
}
