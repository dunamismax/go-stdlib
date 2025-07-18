package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	mathrand "math/rand"
	"time"
)

func init() {
	mathrand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(n.Int64()) + min
}

func RandomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range length {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result[i] = chars[randomIndex.Int64()]
	}
	return string(result)
}

func SecureRandomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range length {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result[i] = chars[randomIndex.Int64()]
	}
	return string(result)
}

func SecureRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	return bytes, err
}

func SecureRandomHex(length int) (string, error) {
	bytes, err := SecureRandomBytes(length)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func RollDice(sides int) int {
	return RandomInt(1, sides)
}
