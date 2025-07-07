package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

func GenerateRandomCode() (string, error) {
	max := big.NewInt(900000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	otp := n.Int64() + 100000

	return fmt.Sprintf("%06d", otp), nil
}

func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	key := hex.EncodeToString(bytes)

	return key, nil
}

func HashAPIKey(prefix, key string) string {
	hashInput := prefix + key
	hash := sha256.Sum256([]byte(hashInput))

	key = hex.EncodeToString(hash[:])

	return key
}
