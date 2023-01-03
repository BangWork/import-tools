package utils

import (
	"math/rand"
	"time"
)

const (
	code         = "0123456789ABCDEFGHIJKLMNOPQRSTUVXWYZabcdefghijklmnopqrstuvxwyz-*"
	numberOffset = 10
)

func RandomNumberString(size int) string {
	return randomString(size, numberOffset, 0)
}

func randomString(size int, max int, seed int64) string {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)
	buffer := make([]byte, size, size)
	for i := 0; i < size; i++ {
		buffer[i] = code[rand.Intn(max)]
	}
	return string(buffer[:size])
}
