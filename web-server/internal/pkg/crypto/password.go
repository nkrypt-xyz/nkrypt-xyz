package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"
)

// HashPassword hashes the given plaintext password using Argon2id and returns
// the hash and salt as strings. The Argon2 parameters are provided explicitly
// so they can be driven from configuration.
func HashPassword(password string, memory, iterations, keyLen uint32, parallelism uint8) (hash string, salt string, err error) {
	if password == "" {
		return "", "", errors.New("password must not be empty")
	}

	saltBytes := make([]byte, 16)
	if _, err = rand.Read(saltBytes); err != nil {
		return "", "", err
	}

	hashBytes := argon2.IDKey([]byte(password), saltBytes, iterations, memory, parallelism, keyLen)

	// Store as base64 to keep compatibility with textual storage.
	return base64.StdEncoding.EncodeToString(hashBytes), base64.StdEncoding.EncodeToString(saltBytes), nil
}

// VerifyPassword verifies that the given plaintext password matches the
// supplied Argon2id hash and salt strings.
func VerifyPassword(password, encodedHash, encodedSalt string, memory, iterations, keyLen uint32, parallelism uint8) (bool, error) {
	if password == "" || encodedHash == "" || encodedSalt == "" {
		return false, errors.New("password, hash, and salt must be non-empty")
	}

	hashBytes, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false, err
	}

	saltBytes, err := base64.StdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return false, err
	}

	computed := argon2.IDKey([]byte(password), saltBytes, iterations, memory, parallelism, keyLen)

	if len(computed) != len(hashBytes) {
		return false, nil
	}

	// Constant-time comparison.
	var diff uint8
	for i := 0; i < len(computed); i++ {
		diff |= computed[i] ^ hashBytes[i]
	}
	return diff == 0, nil
}

