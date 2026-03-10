package util

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

// RandomName returns a random first and last name using gofakeit.
func RandomName() (string, string) {
	return gofakeit.FirstName(), gofakeit.LastName()
}

// RandomBirthdate returns a random birthdate string in YYYY-MM-DD format from 1985-2002.
func RandomBirthdate() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	year := r.Intn(2002-1985+1) + 1985
	month := r.Intn(12) + 1
	day := r.Intn(28) + 1
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}
