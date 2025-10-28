package auth

import (
	"crypto/subtle"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type Argon2Auth struct {
	Argon2Memory   uint32
	Argon2Time     uint32
	Argon2Threads  uint8
	Salt 		   []byte
}

func (auth *Argon2Auth) SetDefaults() error {
	auth.Argon2Memory = 64 * 1024
	auth.Argon2Time = 1
	auth.Argon2Threads = 4
	
	salt, err := generateSalt(16)
	if err != nil {
		return err
	}
	auth.Salt = salt

	return nil
}

func (auth *Argon2Auth) HashPassword(password []byte) []byte {
	return argon2.IDKey([]byte(password), auth.Salt, auth.Argon2Time, auth.Argon2Memory, auth.Argon2Threads, 32)
}

func (auth *Argon2Auth) VerifyPassword(password []byte, expected []byte) bool {
	hashed_password := auth.HashPassword(password)
	return subtle.ConstantTimeCompare(hashed_password, expected) == 1
}

func generateSalt(n int) ([]byte, error) {
	salt := make([]byte, n)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}