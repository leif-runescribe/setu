package utils

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenId() string {
	id := make([]byte, 6)
	for i := range id {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic("Error generating random number")
		}
		id[i] = charset[randomIndex.Int64()]
	}
	return string(id)
}
