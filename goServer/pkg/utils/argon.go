package utils

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

type ArgonConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

var DefaultArgonConfig = ArgonConfig{
	Time:    1,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
}

func Argon2HashPassword(pw string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(pw), salt, DefaultArgonConfig.Time, DefaultArgonConfig.Memory, DefaultArgonConfig.Threads, DefaultArgonConfig.KeyLen)

	return base64.RawStdEncoding.EncodeToString(salt) + "." + base64.RawStdEncoding.EncodeToString(hash), nil
}

func Argon2CheckPassword(pw string, encodedHash string) bool {
	parts := splitOnce(encodedHash, '.')
	if parts == nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	newHash := argon2.IDKey([]byte(pw), salt, DefaultArgonConfig.Time, DefaultArgonConfig.Memory, DefaultArgonConfig.Threads, DefaultArgonConfig.KeyLen)
	return compareBytes(hash, newHash)
}

func splitOnce(s string, sep byte) []string {
	for i := range s {
		if s[i] == sep {
			return []string{s[:i], s[i+1:]}
		}
	}
	return nil
}

func compareBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	result := byte(0)
	for i := range a {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
