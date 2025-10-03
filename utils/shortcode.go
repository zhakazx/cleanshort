package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	shortCodeChars         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
	defaultShortCodeLength = 8
)

func GenerateShortCode(length int) (string, error) {
	if length <= 0 {
		length = defaultShortCodeLength
	}

	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(shortCodeChars))))
		if err != nil {
			return "", err
		}
		result[i] = shortCodeChars[num.Int64()]
	}

	return string(result), nil
}

func IsValidShortCode(shortCode string) bool {
	if len(shortCode) < 4 || len(shortCode) > 32 {
		return false
	}

	for _, char := range shortCode {
		found := false
		for _, allowed := range shortCodeChars {
			if char == allowed {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

var ReservedShortCodes = map[string]bool{
	"api":      true,
	"admin":    true,
	"healthz":  true,
	"readyz":   true,
	"docs":     true,
	"swagger":  true,
	"www":      true,
	"app":      true,
	"auth":     true,
	"login":    true,
	"logout":   true,
	"register": true,
	"signup":   true,
}

func IsReservedShortCode(shortCode string) bool {
	return ReservedShortCodes[shortCode]
}
