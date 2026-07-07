package security

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const testSecret = "abcdefghijklmnopqrstuvwxyz1234567890"

func TestTokenManagerGeneratesAndParsesAccessToken(t *testing.T) {
	fixed := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	manager, err := NewTokenManager(testSecret, 15*time.Minute)
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}
	manager = manager.WithClock(func() time.Time { return fixed })

	userID := uuid.New()
	token, expiresAt, err := manager.GenerateAccessToken(userID)
	if err != nil {
		t.Fatalf("generate access token: %v", err)
	}
	if !expiresAt.Equal(fixed.Add(15 * time.Minute)) {
		t.Fatalf("expiresAt = %s, want %s", expiresAt, fixed.Add(15*time.Minute))
	}

	parsedID, err := manager.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("parse access token: %v", err)
	}
	if parsedID != userID {
		t.Fatalf("parsed user id = %s, want %s", parsedID, userID)
	}
}

func TestTokenManagerRejectsExpiredAccessToken(t *testing.T) {
	fixed := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	manager, err := NewTokenManager(testSecret, time.Minute)
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}
	token, _, err := manager.WithClock(func() time.Time { return fixed }).GenerateAccessToken(uuid.New())
	if err != nil {
		t.Fatalf("generate access token: %v", err)
	}

	_, err = manager.WithClock(func() time.Time { return fixed.Add(2 * time.Minute) }).ParseAccessToken(token)
	if err == nil {
		t.Fatal("expected expired access token to fail")
	}
}

func TestTokenManagerRejectsWrongSecret(t *testing.T) {
	manager, err := NewTokenManager(testSecret, time.Minute)
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}
	token, _, err := manager.GenerateAccessToken(uuid.New())
	if err != nil {
		t.Fatalf("generate access token: %v", err)
	}

	wrongManager, err := NewTokenManager("0123456789abcdefghijklmnopqrstuvwxyz", time.Minute)
	if err != nil {
		t.Fatalf("new wrong token manager: %v", err)
	}
	_, err = wrongManager.ParseAccessToken(token)
	if err == nil {
		t.Fatal("expected wrong signing secret to fail")
	}
}

func TestTokenManagerRejectsMissingSubject(t *testing.T) {
	manager, err := NewTokenManager(testSecret, time.Minute)
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	_, err = manager.ParseAccessToken(token)
	if err == nil {
		t.Fatal("expected missing subject to fail")
	}
}

func TestTokenManagerRejectsMalformedToken(t *testing.T) {
	manager, err := NewTokenManager(testSecret, time.Minute)
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}
	_, err = manager.ParseAccessToken("not-a-jwt")
	if err == nil {
		t.Fatal("expected malformed token to fail")
	}
}
