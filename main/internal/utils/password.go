package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"strings"

	"golang.org/x/crypto/argon2"
)

func VerifyPassword(password, hashedPassword string) bool {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 3 {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		log.Println("Error decoding salt:", err)
		return false
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		log.Println("Error decoding hash:", err)
		return false
	}

	computedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, uint32(len(hash)))

	return base64.RawStdEncoding.EncodeToString(computedHash) == base64.RawStdEncoding.EncodeToString(hash)
}

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return "argon2id$" + base64.RawStdEncoding.EncodeToString(salt) + "$" + base64.RawStdEncoding.EncodeToString(hash), nil
}
