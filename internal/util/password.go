package util

import (
	"math/rand"
	"time"
)

const (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars   = "0123456789"
	specialChars = "!@#$%&*"
	allChars     = lowerChars + upperChars + digitChars + specialChars
)

// GeneratePassword generates a random password of given length (default 14).
// It guarantees at least 1 lowercase, 1 uppercase, 1 digit, and 1 special character.
func GeneratePassword(length int) string {
	if length <= 0 {
		length = 14
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	password := make([]byte, length)
	password[0] = lowerChars[r.Intn(len(lowerChars))]
	password[1] = upperChars[r.Intn(len(upperChars))]
	password[2] = digitChars[r.Intn(len(digitChars))]
	password[3] = specialChars[r.Intn(len(specialChars))]

	for i := 4; i < length; i++ {
		password[i] = allChars[r.Intn(len(allChars))]
	}

	r.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password)
}