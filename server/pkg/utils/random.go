package utils

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"golang.org/x/exp/constraints"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"

// RandomInt generates a random number between min(inclusive) and max(exclusive)
func RandomInt[T constraints.Integer](min, max T) T {
	return min + rand.N(max-min)
}

func RandomItem[T any](items []T) T {
	return items[rand.IntN(len(items))]
}

func RandomItemExcept[T comparable](items []T, except T) T {
	idx := rand.IntN(len(items))
	if items[idx] == except {
		idx = (idx + 1) % len(items)
	}
	return items[idx]
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

// RandomName generates a random string of length 15
func RandomName() string {
	return RandomString(15)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@test.com", RandomName())
}
