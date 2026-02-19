package randstr

import (
	"testing"
)

func TestGenerateIDLength(t *testing.T) {
	id, err := GenerateID(16)
	if err != nil {
		t.Fatalf("GenerateID failed: %v", err)
	}
	if len(id) != 16 {
		t.Errorf("Expected length 16, got %d", len(id))
	}
}

func TestGenerateIDCharset(t *testing.T) {
	id, err := GenerateID(32)
	if err != nil {
		t.Fatalf("GenerateID failed: %v", err)
	}
	for _, c := range id {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			t.Errorf("Invalid character in ID: %c", c)
		}
	}
}

func TestGenerateIDUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		id, err := GenerateID(16)
		if err != nil {
			t.Fatalf("GenerateID failed: %v", err)
		}
		if seen[id] {
			t.Errorf("Duplicate ID generated: %s", id)
		}
		seen[id] = true
	}
}

func TestGenerateAPIKey(t *testing.T) {
	key, err := GenerateAPIKey(128)
	if err != nil {
		t.Fatalf("GenerateAPIKey failed: %v", err)
	}
	if len(key) != 128 {
		t.Errorf("Expected length 128, got %d", len(key))
	}
}
