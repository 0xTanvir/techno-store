package algo

import (
	"math/rand"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// GenerateRandomString generates a random string of the given length.
func GenerateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateRandomInteger generates a random integer between min (inclusive) and max (inclusive).
func GenerateRandomInteger(min, max int) int {
	if min > max {
		min, max = max, min // Swap min and max if they are in the wrong order
	}
	return rand.Intn(max-min+1) + min
}
