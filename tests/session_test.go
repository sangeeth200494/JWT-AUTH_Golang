package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sangeeth200494/JWT-AUTH_Golang/helpers"
	"github.com/sangeeth200494/JWT-AUTH_Golang/login"
	"github.com/sangeeth200494/JWT-AUTH_Golang/middleware"
)

func init() {
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test-secret-key-1234567890-jwt-secret-testing")
	}
}

// TestSessionHandler_Authenticated verifies that an authenticated user receives safe session metadata and a masked token
func TestSessionHandler_Authenticated(t *testing.T) {
	now := time.Now()
	rawToken, err := helpers.CreateToken(42, "secure_user", now, now, now, "active", "admin")
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/auth/session", nil)
	// Provide token via secure HttpOnly cookie
	req.AddCookie(&http.Cookie{
		Name:     "jwt_token",
		Value:    rawToken,
		Path:     "/",
		HttpOnly: true,
	})

	rr := httptest.NewRecorder()

	// Wrap with AuthMiddleware
	handler := middleware.AuthMiddleware(http.HandlerFunc(login.SessionHandler))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, body: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal JSON response: %v", err)
	}

	// 1. Verify authenticated status
	if auth, ok := resp["authenticated"].(bool); !ok || !auth {
		t.Fatalf("expected authenticated=true, got: %v", resp["authenticated"])
	}

	// 2. Verify token type
	if tokenType, ok := resp["token_type"].(string); !ok || tokenType != "Bearer JWT" {
		t.Fatalf("expected token_type='Bearer JWT', got: %v", resp["token_type"])
	}

	// 3. Verify session status
	if sessionStatus, ok := resp["session_status"].(string); !ok || sessionStatus != "active" {
		t.Fatalf("expected session_status='active', got: %v", resp["session_status"])
	}

	// 4. Verify token status
	if tokenStatus, ok := resp["token_status"].(string); !ok || tokenStatus != "Valid" {
		t.Fatalf("expected token_status='Valid', got: %v", resp["token_status"])
	}

	// 5. Verify masked token
	maskedToken, ok := resp["masked_token"].(string)
	if !ok || maskedToken == "" {
		t.Fatalf("expected masked_token to be non-empty string")
	}
	if !strings.Contains(maskedToken, "••••••••••••••••") {
		t.Fatalf("expected masked_token to contain masking bullets, got: %s", maskedToken)
	}

	// 6. CRITICAL SECURITY CHECK: The raw token must NEVER appear anywhere in the response!
	if strings.Contains(rr.Body.String(), rawToken) {
		t.Fatalf("CRITICAL SECURITY VIOLATION: Raw JWT token was leaked in session response!")
	}
}

// TestSessionHandler_Unauthenticated verifies that unauthenticated requests are rejected
func TestSessionHandler_Unauthenticated(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/auth/session", nil)
	rr := httptest.NewRecorder()

	handler := middleware.AuthMiddleware(http.HandlerFunc(login.SessionHandler))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got: %d", rr.Code)
	}
}

// TestLogoutHandler verifies that the logout endpoint clears the HttpOnly cookie
func TestLogoutHandler(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	rr := httptest.NewRecorder()

	login.LogoutHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", rr.Code)
	}

	cookies := rr.Result().Cookies()
	var jwtCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "jwt_token" {
			jwtCookie = c
			break
		}
	}

	if jwtCookie == nil {
		t.Fatalf("expected jwt_token cookie to be present in response to clear it")
	}
	if jwtCookie.MaxAge != -1 {
		t.Fatalf("expected cookie MaxAge to be -1, got: %d", jwtCookie.MaxAge)
	}
	if jwtCookie.Value != "" {
		t.Fatalf("expected cookie value to be empty, got: %s", jwtCookie.Value)
	}
}
