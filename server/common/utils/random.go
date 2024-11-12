package utils

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"

// RandomInt generates a random number between min and max
func RandomInt(min, max int) int {
	if min <= 0 {
		min = 1
	}
	if max <= 0 {
		max = 1
	}
	return min + rand.IntN(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabets)

	for i := 0; i < n; i++ {
		c := alphabets[rand.IntN(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomName generates a small random name
func RandomName() string {
	return RandomString(8)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomName())
}
