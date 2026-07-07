package service

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/WassimBannout/gatherops/internal/domain"
	"github.com/WassimBannout/gatherops/internal/repository"
	"github.com/WassimBannout/gatherops/internal/security"
	"github.com/google/uuid"
)

const minPasswordLength = 8

type AuthService struct {
	users         repository.UserRepository
	refreshTokens repository.RefreshTokenRepository
	passwords     security.PasswordHasher
	tokens        security.TokenManager
	refreshTTL    time.Duration
	now           func() time.Time
}

type AuthServiceConfig struct {
	Users         repository.UserRepository
	RefreshTokens repository.RefreshTokenRepository
	Passwords     security.PasswordHasher
	Tokens        security.TokenManager
	RefreshTTL    time.Duration
	Now           func() time.Time
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthSession struct {
	User                  UserProfile
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

type UserProfile struct {
	ID        uuid.UUID
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAuthService(cfg AuthServiceConfig) (*AuthService, error) {
	if cfg.Users == nil {
		return nil, errors.New("user repository is required")
	}
	if cfg.RefreshTokens == nil {
		return nil, errors.New("refresh token repository is required")
	}
	if cfg.RefreshTTL <= 0 {
		return nil, errors.New("refresh token ttl must be positive")
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}

	return &AuthService{
		users:         cfg.Users,
		refreshTokens: cfg.RefreshTokens,
		passwords:     cfg.Passwords,
		tokens:        cfg.Tokens,
		refreshTTL:    cfg.RefreshTTL,
		now:           cfg.Now,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (AuthSession, error) {
	name := strings.TrimSpace(input.Name)
	email := normalizeEmail(input.Email)
	details := validateNameEmailPassword(name, email, input.Password)
	if len(details) > 0 {
		return AuthSession{}, NewError(ErrorCodeValidationFailed, "One or more fields are invalid", details)
	}

	passwordHash, err := s.passwords.Hash(input.Password)
	if err != nil {
		return AuthSession{}, err
	}

	user, err := s.users.Create(ctx, domain.User{
		Email:        email,
		Name:         name,
		PasswordHash: passwordHash,
	})
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return AuthSession{}, NewError(ErrorCodeEmailAlreadyRegistered, "Email is already registered", map[string]any{"email": "already registered"})
		}
		return AuthSession{}, err
	}

	return s.issueSession(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (AuthSession, error) {
	email := normalizeEmail(input.Email)
	if !isValidEmail(email) || strings.TrimSpace(input.Password) == "" {
		return AuthSession{}, invalidCredentials()
	}

	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return AuthSession{}, invalidCredentials()
		}
		return AuthSession{}, err
	}

	if !s.passwords.Compare(user.PasswordHash, input.Password) {
		return AuthSession{}, invalidCredentials()
	}

	return s.issueSession(ctx, user)
}

func (s *AuthService) Refresh(ctx context.Context, rawRefreshToken string) (AuthSession, error) {
	token := strings.TrimSpace(rawRefreshToken)
	if token == "" {
		return AuthSession{}, invalidRefreshToken()
	}

	stored, err := s.refreshTokens.FindByHash(ctx, security.HashRefreshToken(token))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return AuthSession{}, invalidRefreshToken()
		}
		return AuthSession{}, err
	}
	if stored.RevokedAt != nil || !stored.ExpiresAt.After(s.now().UTC()) {
		return AuthSession{}, invalidRefreshToken()
	}

	user, err := s.users.FindByID(ctx, stored.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return AuthSession{}, invalidRefreshToken()
		}
		return AuthSession{}, err
	}

	if err := s.refreshTokens.Revoke(ctx, stored.ID); err != nil {
		return AuthSession{}, err
	}

	return s.issueSession(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID, rawRefreshToken string) error {
	if userID == uuid.Nil {
		return NewError(ErrorCodeUnauthorized, "Authentication required", nil)
	}
	token := strings.TrimSpace(rawRefreshToken)
	if token == "" {
		return invalidRefreshToken()
	}

	stored, err := s.refreshTokens.FindByHash(ctx, security.HashRefreshToken(token))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return invalidRefreshToken()
		}
		return err
	}
	if stored.UserID != userID {
		return NewError(ErrorCodeForbidden, "Refresh token does not belong to the authenticated user", nil)
	}
	if stored.RevokedAt != nil {
		return nil
	}

	return s.refreshTokens.Revoke(ctx, stored.ID)
}

func (s *AuthService) Me(ctx context.Context, userID uuid.UUID) (UserProfile, error) {
	if userID == uuid.Nil {
		return UserProfile{}, NewError(ErrorCodeUnauthorized, "Authentication required", nil)
	}

	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return UserProfile{}, NewError(ErrorCodeUnauthorized, "Authentication required", nil)
		}
		return UserProfile{}, err
	}

	return profileFromUser(user), nil
}

func (s *AuthService) issueSession(ctx context.Context, user domain.User) (AuthSession, error) {
	accessToken, accessExpiresAt, err := s.tokens.GenerateAccessToken(user.ID)
	if err != nil {
		return AuthSession{}, err
	}

	rawRefreshToken, err := security.GenerateRefreshToken()
	if err != nil {
		return AuthSession{}, err
	}
	refreshExpiresAt := s.now().UTC().Add(s.refreshTTL)

	if _, err := s.refreshTokens.Create(ctx, domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: security.HashRefreshToken(rawRefreshToken),
		ExpiresAt: refreshExpiresAt,
	}); err != nil {
		return AuthSession{}, err
	}

	return AuthSession{
		User:                  profileFromUser(user),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshToken:          rawRefreshToken,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}, nil
}

func profileFromUser(user domain.User) UserProfile {
	return UserProfile{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func validateNameEmailPassword(name, email, password string) map[string]any {
	details := map[string]any{}
	if name == "" {
		details["name"] = "is required"
	}
	if !isValidEmail(email) {
		details["email"] = "must be a valid email address"
	}
	if len(password) < minPasswordLength {
		details["password"] = "must be at least 8 characters"
	}
	if len(details) == 0 {
		return nil
	}
	return details
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func isValidEmail(email string) bool {
	if email == "" || strings.ContainsAny(email, " \t\n\r") {
		return false
	}
	addr, err := mail.ParseAddress(email)
	return err == nil && addr.Address == email
}

func invalidCredentials() error {
	return NewError(ErrorCodeInvalidCredentials, "Invalid email or password", nil)
}

func invalidRefreshToken() error {
	return NewError(ErrorCodeInvalidRefreshToken, "Refresh token is invalid or expired", nil)
}
