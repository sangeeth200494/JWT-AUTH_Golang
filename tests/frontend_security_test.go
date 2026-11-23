package tests

import (
	"os"
	"strings"
	"testing"
)

// TestFrontend_NoTokenStorageLeaks verifies that frontend JavaScript does NOT store JWT tokens in localStorage or sessionStorage
func TestFrontend_NoTokenStorageLeaks(t *testing.T) {
	appJs, err := os.ReadFile("../static/app.js")
	if err != nil {
		t.Fatalf("failed to read static/app.js: %v", err)
	}
	content := string(appJs)

	// Verify no localStorage storage of jwt
	if strings.Contains(content, "localStorage.setItem") {
		t.Fatalf("SECURITY VIOLATION: static/app.js must not store credentials in localStorage!")
	}

	// Verify no sessionStorage storage of jwt
	if strings.Contains(content, "sessionStorage.setItem") {
		t.Fatalf("SECURITY VIOLATION: static/app.js must not store credentials in sessionStorage!")
	}

	// Verify interaction with safe session endpoint
	if !strings.Contains(content, "/api/v1/auth/session") {
		t.Fatalf("expected static/app.js to retrieve dynamic metadata from /api/v1/auth/session")
	}

	// Verify credentials: 'same-origin' is used for HttpOnly cookie transmission
	if !strings.Contains(content, "credentials: 'same-origin'") {
		t.Fatalf("expected static/app.js to use credentials: 'same-origin' for cookie transport")
	}
}

// TestFrontend_MaskedTokenComponent verifies that the UI contains the dedicated Authentication Details component
func TestFrontend_MaskedTokenComponent(t *testing.T) {
	html, err := os.ReadFile("../static/index.html")
	if err != nil {
		t.Fatalf("failed to read static/index.html: %v", err)
	}
	content := string(html)

	// Verify header
	if !strings.Contains(content, "AUTHENTICATION DETAILS") {
		t.Fatalf("expected static/index.html to contain 'AUTHENTICATION DETAILS' section")
	}

	// Verify masked token container
	if !strings.Contains(content, "auth-masked-token") {
		t.Fatalf("expected static/index.html to contain 'auth-masked-token' element")
	}

	// Verify raw token container has been removed
	if strings.Contains(content, "dash-raw-token") {
		t.Fatalf("SECURITY VIOLATION: static/index.html still contains 'dash-raw-token' raw display!")
	}
}
