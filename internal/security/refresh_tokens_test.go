package security

import "testing"

func TestRefreshTokenGenerationAndHashing(t *testing.T) {
	token, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("generate refresh token: %v", err)
	}
	if token == "" {
		t.Fatal("expected refresh token")
	}
	hash := HashRefreshToken(token)
	if hash == token {
		t.Fatal("refresh token hash must not equal raw token")
	}
	if len(hash) != 64 {
		t.Fatalf("hash length = %d, want 64", len(hash))
	}
	if HashRefreshToken(token) != hash {
		t.Fatal("refresh token hashing must be deterministic")
	}
}
