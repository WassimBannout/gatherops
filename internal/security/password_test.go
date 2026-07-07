package security

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestPasswordHasherHashesAndCompares(t *testing.T) {
	hasher := NewPasswordHasher(bcrypt.MinCost)
	hash, err := hasher.Hash("correct horse battery staple")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if hash == "correct horse battery staple" {
		t.Fatal("password hash must not equal raw password")
	}
	if !hasher.Compare(hash, "correct horse battery staple") {
		t.Fatal("expected password comparison to succeed")
	}
	if hasher.Compare(hash, "wrong password") {
		t.Fatal("expected wrong password comparison to fail")
	}
}
