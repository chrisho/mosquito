package helper

import (
	"time"
	"math/rand"
	"strings"
	"strconv"
)

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	rs := make([]string, length)

	for start := 0; start < length; start++ {
		t := rand.Intn(2)
		if t == 0 {
			// A-Z
			rs = append(rs, string(rand.Intn(26)+65))
		} else {
			// a-z
			rs = append(rs, string(rand.Intn(26)+97))
		}
	}
	return strings.Join(rs, "")
}

func RandomInt(n int) int {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	return random.Intn(n)
}

func RandomIntString(n int) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	return strconv.Itoa(random.Intn(n))
}

// 长度为 length 的整形字符串(允许 0 开始)
func RandomIntStringLength(length int) string {
	result := ""
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		result += strconv.Itoa(random.Intn(10))
	}

	return result
}

func RandomBetweenInt(min, max int) int {
	if min >= max {
		return min
	}

	return RandomInt(max - min + 1) + min
}
