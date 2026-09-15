package services

import (
	"crypto/rand"
	"encoding/base64"
)

func shortCodeGenerator() (string, error) {
	length := 6
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
