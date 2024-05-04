package random

import (
	"fmt"
	"math/rand"
	"strings"
)

const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(alpha[rand.Intn(len(alpha))])
	}
	return sb.String()
}

func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandomBool() bool {
	return rand.Intn(2) == 1
}

func RandomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomPassword() string {
	return RandomString(8)
}

func RandomName() string {
	return RandomString(6)
}

func RandomBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rand.Intn(256))
	}
	return b
}
