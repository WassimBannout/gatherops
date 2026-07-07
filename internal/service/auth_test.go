package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/WassimBannout/gatherops/internal/domain"
	"github.com/WassimBannout/gatherops/internal/repository"
	"github.com/WassimBannout/gatherops/internal/security"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthServiceRegisterNormalizesEmailAndStoresOnlyHashes(t *testing.T) {
	svc, users, refreshTokens := newTestAuthService(t)

	session, err := svc.Register(context.Background(), RegisterInput{
		Name:     "  Ada Lovelace  ",
		Email:    " ADA@Example.COM ",
		Password: "correct-password",
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	if session.User.Email != "ada@example.com" {
		t.Fatalf("email = %q, want normalized", session.User.Email)
	}
	stored := users.byEmail["ada@example.com"]
	if stored.PasswordHash == "correct-password" {
		t.Fatal("password hash must not store raw password")
	}
	if !security.NewPasswordHasher(bcrypt.MinCost).Compare(stored.PasswordHash, "correct-password") {
		t.Fatal("stored password hash should verify raw password")
	}
	if len(refreshTokens.byHash) != 1 {
		t.Fatalf("stored refresh tokens = %d, want 1", len(refreshTokens.byHash))
	}
	if _, ok := refreshTokens.byHash[session.RefreshToken]; ok {
		t.Fatal("raw refresh token must not be stored as repository key")
	}
	if _, ok := refreshTokens.byHash[security.HashRefreshToken(session.RefreshToken)]; !ok {
		t.Fatal("hashed refresh token should be stored")
	}
}

func TestAuthServiceRegisterRejectsValidationErrors(t *testing.T) {
	svc, _, _ := newTestAuthService(t)

	_, err := svc.Register(context.Background(), RegisterInput{
		Name:     "",
		Email:    "bad-email",
		Password: "short",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	serviceErr := requireServiceError(t, err)
	if serviceErr.Code != ErrorCodeValidationFailed {
		t.Fatalf("code = %s, want validation_failed", serviceErr.Code)
	}
	if len(serviceErr.Details) != 3 {
		t.Fatalf("details = %#v, want 3 validation details", serviceErr.Details)
	}
}

func TestAuthServiceRegisterDuplicateEmailReturnsConflict(t *testing.T) {
	svc, _, _ := newTestAuthService(t)

	input := RegisterInput{Name: "Ada", Email: "ada@example.com", Password: "correct-password"}
	if _, err := svc.Register(context.Background(), input); err != nil {
		t.Fatalf("first register: %v", err)
	}
	_, err := svc.Register(context.Background(), input)
	if err == nil {
		t.Fatal("expected duplicate email to fail")
	}
	serviceErr := requireServiceError(t, err)
	if serviceErr.Code != ErrorCodeEmailAlreadyRegistered {
		t.Fatalf("code = %s, want email_already_registered", serviceErr.Code)
	}
}

func TestAuthServiceLoginRejectsInvalidPasswordWithGenericError(t *testing.T) {
	svc, _, _ := newTestAuthService(t)
	if _, err := svc.Register(context.Background(), RegisterInput{Name: "Ada", Email: "ada@example.com", Password: "correct-password"}); err != nil {
		t.Fatalf("register: %v", err)
	}

	_, err := svc.Login(context.Background(), LoginInput{Email: "ada@example.com", Password: "wrong-password"})
	if err == nil {
		t.Fatal("expected login to fail")
	}
	serviceErr := requireServiceError(t, err)
	if serviceErr.Code != ErrorCodeInvalidCredentials {
		t.Fatalf("code = %s, want invalid_credentials", serviceErr.Code)
	}
	if serviceErr.Message != "Invalid email or password" {
		t.Fatalf("message = %q, want generic login failure", serviceErr.Message)
	}
}

func TestAuthServiceRefreshRotatesAndRejectsReusedRefreshToken(t *testing.T) {
	svc, _, _ := newTestAuthService(t)
	session, err := svc.Register(context.Background(), RegisterInput{Name: "Ada", Email: "ada@example.com", Password: "correct-password"})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	rotated, err := svc.Refresh(context.Background(), session.RefreshToken)
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if rotated.RefreshToken == session.RefreshToken {
		t.Fatal("refresh should rotate refresh token")
	}

	_, err = svc.Refresh(context.Background(), session.RefreshToken)
	if err == nil {
		t.Fatal("expected reused refresh token to fail")
	}
	serviceErr := requireServiceError(t, err)
	if serviceErr.Code != ErrorCodeInvalidRefreshToken {
		t.Fatalf("code = %s, want invalid_refresh_token", serviceErr.Code)
	}
}

func TestAuthServiceLogoutRevokesRefreshToken(t *testing.T) {
	svc, _, _ := newTestAuthService(t)
	session, err := svc.Register(context.Background(), RegisterInput{Name: "Ada", Email: "ada@example.com", Password: "correct-password"})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	if err := svc.Logout(context.Background(), session.User.ID, session.RefreshToken); err != nil {
		t.Fatalf("logout: %v", err)
	}
	_, err = svc.Refresh(context.Background(), session.RefreshToken)
	if err == nil {
		t.Fatal("expected refresh after logout to fail")
	}
	serviceErr := requireServiceError(t, err)
	if serviceErr.Code != ErrorCodeInvalidRefreshToken {
		t.Fatalf("code = %s, want invalid_refresh_token", serviceErr.Code)
	}
}

func TestAuthServiceMeReturnsProfileWithoutPasswordHash(t *testing.T) {
	svc, _, _ := newTestAuthService(t)
	session, err := svc.Register(context.Background(), RegisterInput{Name: "Ada", Email: "ada@example.com", Password: "correct-password"})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	profile, err := svc.Me(context.Background(), session.User.ID)
	if err != nil {
		t.Fatalf("me: %v", err)
	}
	if profile.Email != "ada@example.com" || profile.Name != "Ada" {
		t.Fatalf("profile = %#v", profile)
	}
}

func newTestAuthService(t *testing.T) (*AuthService, *memoryUserRepository, *memoryRefreshTokenRepository) {
	t.Helper()

	fixed := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	users := newMemoryUserRepository(fixed)
	refreshTokens := newMemoryRefreshTokenRepository(fixed)
	tokens, err := security.NewTokenManager("abcdefghijklmnopqrstuvwxyz1234567890", 15*time.Minute)
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}
	tokens = tokens.WithClock(func() time.Time { return fixed })

	svc, err := NewAuthService(AuthServiceConfig{
		Users:         users,
		RefreshTokens: refreshTokens,
		Passwords:     security.NewPasswordHasher(bcrypt.MinCost),
		Tokens:        tokens,
		RefreshTTL:    24 * time.Hour,
		Now:           func() time.Time { return fixed },
	})
	if err != nil {
		t.Fatalf("new auth service: %v", err)
	}

	return svc, users, refreshTokens
}

func requireServiceError(t *testing.T, err error) *Error {
	t.Helper()

	var serviceErr *Error
	if !errors.As(err, &serviceErr) {
		t.Fatalf("error = %T %v, want service error", err, err)
	}
	return serviceErr
}

type memoryUserRepository struct {
	now     time.Time
	byID    map[uuid.UUID]domain.User
	byEmail map[string]domain.User
}

func newMemoryUserRepository(now time.Time) *memoryUserRepository {
	return &memoryUserRepository{
		now:     now,
		byID:    map[uuid.UUID]domain.User{},
		byEmail: map[string]domain.User{},
	}
}

func (r *memoryUserRepository) Create(_ context.Context, user domain.User) (domain.User, error) {
	if _, ok := r.byEmail[user.Email]; ok {
		return domain.User{}, repository.ErrConflict
	}
	user.ID = uuid.New()
	user.CreatedAt = r.now
	user.UpdatedAt = r.now
	r.byID[user.ID] = user
	r.byEmail[user.Email] = user
	return user, nil
}

func (r *memoryUserRepository) FindByID(_ context.Context, id uuid.UUID) (domain.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return domain.User{}, repository.ErrNotFound
	}
	return user, nil
}

func (r *memoryUserRepository) FindByEmail(_ context.Context, email string) (domain.User, error) {
	user, ok := r.byEmail[email]
	if !ok {
		return domain.User{}, repository.ErrNotFound
	}
	return user, nil
}

type memoryRefreshTokenRepository struct {
	now    time.Time
	byID   map[uuid.UUID]domain.RefreshToken
	byHash map[string]domain.RefreshToken
}

func newMemoryRefreshTokenRepository(now time.Time) *memoryRefreshTokenRepository {
	return &memoryRefreshTokenRepository{
		now:    now,
		byID:   map[uuid.UUID]domain.RefreshToken{},
		byHash: map[string]domain.RefreshToken{},
	}
}

func (r *memoryRefreshTokenRepository) Create(_ context.Context, token domain.RefreshToken) (domain.RefreshToken, error) {
	if _, ok := r.byHash[token.TokenHash]; ok {
		return domain.RefreshToken{}, repository.ErrConflict
	}
	token.ID = uuid.New()
	token.CreatedAt = r.now
	r.byID[token.ID] = token
	r.byHash[token.TokenHash] = token
	return token, nil
}

func (r *memoryRefreshTokenRepository) FindByHash(_ context.Context, tokenHash string) (domain.RefreshToken, error) {
	token, ok := r.byHash[tokenHash]
	if !ok {
		return domain.RefreshToken{}, repository.ErrNotFound
	}
	return token, nil
}

func (r *memoryRefreshTokenRepository) Revoke(_ context.Context, id uuid.UUID) error {
	token, ok := r.byID[id]
	if !ok {
		return repository.ErrNotFound
	}
	revokedAt := r.now
	token.RevokedAt = &revokedAt
	r.byID[id] = token
	r.byHash[token.TokenHash] = token
	return nil
}
