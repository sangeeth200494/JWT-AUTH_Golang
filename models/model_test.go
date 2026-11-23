package models

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "SecurePassword123!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error while hashing password, got: %v", err)
	}
	if hash == "" {
		t.Fatal("expected hash to be non-empty")
	}
	if hash == password {
		t.Fatal("expected hash to not equal plaintext password")
	}

	// Verify that hashing the same password twice yields different salts
	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error on second hash, got: %v", err)
	}
	if hash == hash2 {
		t.Fatal("expected different hashes due to salt")
	}
}

func TestValidatePasswords(t *testing.T) {
	password := "MySecretPassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	// Correct password
	if err := ValidatePasswords(password, hash); err != nil {
		t.Fatalf("expected password to validate successfully, got: %v", err)
	}

	// Incorrect password
	if err := ValidatePasswords("WrongPassword", hash); err == nil {
		t.Fatal("expected validation to fail for incorrect password, but it succeeded")
	}
}

func TestAPIResponseError(t *testing.T) {
	resp := &APIResponse{
		Code:    400,
		Message: "Bad Request",
		Details: "Invalid field",
	}
	expected := "Error400:Bad Request"
	if resp.Error() != expected {
		t.Fatalf("expected error string %q, got %q", expected, resp.Error())
	}
}
