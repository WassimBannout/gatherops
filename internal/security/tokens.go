package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidAccessToken = errors.New("invalid access token")

type TokenManager struct {
	secret    []byte
	accessTTL time.Duration
	now       func() time.Time
}

func NewTokenManager(secret string, accessTTL time.Duration) (TokenManager, error) {
	if secret == "" {
		return TokenManager{}, errors.New("JWT access secret is required")
	}
	if accessTTL <= 0 {
		return TokenManager{}, errors.New("access token ttl must be positive")
	}

	return TokenManager{
		secret:    []byte(secret),
		accessTTL: accessTTL,
		now:       time.Now,
	}, nil
}

func (m TokenManager) WithClock(now func() time.Time) TokenManager {
	m.now = now
	return m
}

func (m TokenManager) GenerateAccessToken(userID uuid.UUID) (string, time.Time, error) {
	if userID == uuid.Nil {
		return "", time.Time{}, errors.New("user id is required")
	}

	issuedAt := m.now().UTC()
	expiresAt := issuedAt.Add(m.accessTTL)
	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		ID:        uuid.NewString(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign access token: %w", err)
	}

	return token, expiresAt, nil
}

func (m TokenManager) ParseAccessToken(raw string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	parser := jwt.NewParser(
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithTimeFunc(func() time.Time { return m.now().UTC() }),
	)

	token, err := parser.ParseWithClaims(raw, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidAccessToken
		}
		return m.secret, nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, ErrInvalidAccessToken
	}
	if claims.Subject == "" {
		return uuid.Nil, ErrInvalidAccessToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil || userID == uuid.Nil {
		return uuid.Nil, ErrInvalidAccessToken
	}

	return userID, nil
}
