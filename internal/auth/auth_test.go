package auth

import (
	"testing"
)

func TestArgon2Auth_HashAndVerifyPassword(t *testing.T) {
	auth := &Argon2Auth{}
	err := auth.SetDefaults()
	if err != nil {
		t.Fatalf("SetDefaults() error: %v", err)
	}

	password := []byte("securepassword")
	hashed_password := auth.HashPassword(password)

	if !auth.VerifyPassword(password, hashed_password) {
		t.Errorf("VerifyPassword() failed: password won't match its hash")
	}

	wrong_password := []byte("wrongpassword")

	if auth.VerifyPassword(wrong_password, hashed_password) {
		t.Errorf("VerifyPassword() failed: expected password not to match")
	}
}