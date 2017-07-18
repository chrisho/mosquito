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
