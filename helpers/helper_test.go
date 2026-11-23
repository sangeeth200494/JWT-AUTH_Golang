package helpers

import (
	"os"
	"strings"
	"testing"
	"time"
)

func init() {
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test-secret-key-1234567890-jwt-secret-testing")
	}
}

func TestCreateAndVerifyToken(t *testing.T) {
	userID := uint64(101)
	username := "testuser"
	now := time.Now()
	createdAt := now
	updatedAt := now
	lastLogin := now
	status := "active"
	role := "user"

	tokenString, err := CreateToken(userID, username, createdAt, updatedAt, lastLogin, status, role)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}
	if tokenString == "" {
		t.Fatal("expected non-empty token string")
	}

	// Verify token with "Bearer " prefix
	if err := VerifyToken("Bearer " + tokenString); err != nil {
		t.Fatalf("failed to verify token with Bearer prefix: %v", err)
	}

	// Verify token without "Bearer " prefix
	if err := VerifyToken(tokenString); err != nil {
		t.Fatalf("failed to verify token without Bearer prefix: %v", err)
	}
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	invalidToken := "invalid.token.string"
	if err := VerifyToken(invalidToken); err == nil {
		t.Fatal("expected error verifying invalid token, got nil")
	}
}

func TestExtractUsernameFromToken(t *testing.T) {
	userID := uint64(42)
	expectedUsername := "alice"
	now := time.Now()

	tokenString, err := CreateToken(userID, expectedUsername, now, now, now, "active", "admin")
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	extracted, err := ExtractUsernameFromToken("Bearer " + tokenString)
	if err != nil {
		t.Fatalf("failed to extract username: %v", err)
	}

	if extracted != expectedUsername {
		t.Fatalf("expected extracted username %q, got %q", expectedUsername, extracted)
	}
}

func TestMaskToken(t *testing.T) {
	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU0NTU1NDB9.signaturePartHere"
	masked := MaskToken(token)

	// Verify prefix
	if !strings.HasPrefix(masked, "eyJhbGciOiJIUzUx...") {
		t.Fatalf("expected masked token to start with prefix, got: %s", masked)
	}

	// Verify masking bullets
	if !strings.Contains(masked, "••••••••••••••••") {
		t.Fatalf("expected masked token to contain bullets, got: %s", masked)
	}

	// Verify full token is NEVER exposed
	if strings.Contains(masked, "signaturePartHere") {
		t.Fatalf("masked token must never expose raw signature or full token!")
	}

	// Test with Bearer prefix
	maskedBearer := MaskToken("Bearer " + token)
	if maskedBearer != masked {
		t.Fatalf("expected MaskToken with Bearer prefix to match, got: %s vs %s", maskedBearer, masked)
	}

	// Test empty or short
	if MaskToken("") != "••••••••••••••••" {
		t.Fatalf("expected empty token to return bullets")
	}
}

func TestExtractClaimsFromToken(t *testing.T) {
	userID := uint64(999)
	username := "bob"
	now := time.Now()

	tokenString, err := CreateToken(userID, username, now, now, now, "active", "admin")
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	claims, err := ExtractClaimsFromToken(tokenString)
	if err != nil {
		t.Fatalf("failed to extract claims: %v", err)
	}

	if claims["username"] != username {
		t.Fatalf("expected username claim %s, got %v", username, claims["username"])
	}

	if exp, ok := claims["exp"].(float64); !ok || exp <= 0 {
		t.Fatalf("expected valid exp claim, got %v", claims["exp"])
	}
}
