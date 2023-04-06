package util

import (
	"math/rand"
	"strings"
)

var (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

// generate random number
func GenRandomNum(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// generate random string
func GenRandomString(n int) string {
	var str strings.Builder
	l := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(l)]
		str.WriteByte(c)
	}

	return str.String()
}

// generate random owner
func GenRandomOwner() string {
	return GenRandomString(8)
}

// generate random money
func GenRandomMoney() int64 {
	return GenRandomNum(0, 1000000)
}

// generate random currency
func GenRandomCurrency() string {
	cr := []string{"EUR", "USD", "IDR"}
	n := len(cr)
	return cr[rand.Intn(n)]
}
