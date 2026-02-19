package crypto

import (
	"testing"
)

// Default Argon2 parameters for testing
const (
	testMemory      = 64 * 1024
	testIterations  = 3
	testParallelism = 2
	testKeyLength   = 32
)

func TestHashAndVerify(t *testing.T) {
	password := "TestPassword123!"
	hash, salt, err := HashPassword(password, testMemory, testIterations, testKeyLength, testParallelism)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	match, err := VerifyPassword(password, hash, salt, testMemory, testIterations, testKeyLength, testParallelism)
	if err != nil {
		t.Fatalf("VerifyPassword failed: %v", err)
	}
	if !match {
		t.Error("VerifyPassword failed for correct password")
	}
}

func TestVerifyWrongPassword(t *testing.T) {
	password := "TestPassword123!"
	wrongPassword := "WrongPassword456!"
	hash, salt, err := HashPassword(password, testMemory, testIterations, testKeyLength, testParallelism)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	match, err := VerifyPassword(wrongPassword, hash, salt, testMemory, testIterations, testKeyLength, testParallelism)
	if err != nil {
		t.Fatalf("VerifyPassword failed: %v", err)
	}
	if match {
		t.Error("VerifyPassword succeeded for wrong password")
	}
}

func TestDifferentSalts(t *testing.T) {
	password := "TestPassword123!"
	hash1, salt1, err := HashPassword(password, testMemory, testIterations, testKeyLength, testParallelism)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	hash2, salt2, err := HashPassword(password, testMemory, testIterations, testKeyLength, testParallelism)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if salt1 == salt2 {
		t.Error("Two hashes of same password have identical salts")
	}
	if hash1 == hash2 {
		t.Error("Two hashes of same password are identical (should differ due to salt)")
	}
}

func TestHashLength(t *testing.T) {
	password := "TestPassword123!"
	hash, salt, err := HashPassword(password, testMemory, testIterations, testKeyLength, testParallelism)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if len(hash) == 0 {
		t.Error("Hash is empty")
	}
	if len(salt) == 0 {
		t.Error("Salt is empty")
	}
}
