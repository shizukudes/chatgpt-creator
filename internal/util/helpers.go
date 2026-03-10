package util

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandStr generates a random alphanumeric string of given length.
func RandStr(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = alphanumeric[r.Intn(len(alphanumeric))]
	}
	return string(b)
}
// GenerateUUID generates a random UUID.
func GenerateUUID() string {
	return uuid.New().String()
}