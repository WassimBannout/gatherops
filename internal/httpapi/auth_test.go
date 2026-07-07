package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WassimBannout/gatherops/internal/domain"
	"github.com/WassimBannout/gatherops/internal/repository"
	"github.com/WassimBannout/gatherops/internal/security"
	"github.com/WassimBannout/gatherops/internal/service"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterReturnsSessionWithoutPasswordHash(t *testing.T) {
	router, _ := newAuthTestRouter(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", jsonBody(`{
		"name":"Ada Lovelace",
		"email":"ADA@example.com",
		"password":"correct-password"
	}`))
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	body := rec.Body.String()
	if strings.Contains(body, "passwordHash") || strings.Contains(body, "correct-password") {
		t.Fatalf("response leaked password material: %s", body)
	}

	var response dataResponse[authSessionResponse]
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Data.User.Email != "ada@example.com" {
		t.Fatalf("email = %q, want normalized", response.Data.User.Email)
	}
	if response.Data.AccessToken == "" || response.Data.RefreshToken == "" {
		t.Fatal("expected access and refresh tokens")
	}
}

func TestRegisterDuplicateEmailReturnsConflict(t *testing.T) {
	router, _ := newAuthTestRouter(t)
	payload := `{"name":"Ada","email":"ada@example.com","password":"correct-password"}`

	first := httptest.NewRecorder()
	router.ServeHTTP(first, httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", jsonBody(payload)))
	if first.Code != http.StatusCreated {
		t.Fatalf("first status = %d, want 201; body: %s", first.Code, first.Body.String())
	}

	second := httptest.NewRecorder()
	router.ServeHTTP(second, httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", jsonBody(payload)))
	if second.Code != http.StatusConflict {
		t.Fatalf("second status = %d, want 409; body: %s", second.Code, second.Body.String())
	}

	var response ErrorResponse
	if err := json.Unmarshal(second.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if response.Error.Code != "email_already_registered" {
		t.Fatalf("error code = %q, want email_already_registered", response.Error.Code)
	}
}

func TestRegisterValidationFailureReturnsDetails(t *testing.T) {
	router, _ := newAuthTestRouter(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", jsonBody(`{
		"name":"",
		"email":"bad-email",
		"password":"short"
	}`))
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body: %s", rec.Code, rec.Body.String())
	}
	var response ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Error.Code != "validation_failed" {
		t.Fatalf("code = %q, want validation_failed", response.Error.Code)
	}
	if len(response.Error.Details) != 3 {
		t.Fatalf("details = %#v, want 3 fields", response.Error.Details)
	}
}

func TestProtectedEndpointRejectsMalformedToken(t *testing.T) {
	router, _ := newAuthTestRouter(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req.Header.Set("Authorization", "Bearer not-a-jwt")
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401; body: %s", rec.Code, rec.Body.String())
	}
	var response ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Error.Code != "invalid_access_token" {
		t.Fatalf("code = %q, want invalid_access_token", response.Error.Code)
	}
}

func TestLogoutRevokesRefreshToken(t *testing.T) {
	router, _ := newAuthTestRouter(t)

	registerRec := httptest.NewRecorder()
	router.ServeHTTP(registerRec, httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", jsonBody(`{
		"name":"Ada",
		"email":"ada@example.com",
		"password":"correct-password"
	}`)))
	if registerRec.Code != http.StatusCreated {
		t.Fatalf("register status = %d, want 201; body: %s", registerRec.Code, registerRec.Body.String())
	}

	var session dataResponse[authSessionResponse]
	if err := json.Unmarshal(registerRec.Body.Bytes(), &session); err != nil {
		t.Fatalf("decode register response: %v", err)
	}

	logoutRec := httptest.NewRecorder()
	logoutReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", jsonBody(`{"refreshToken":"`+session.Data.RefreshToken+`"}`))
	logoutReq.Header.Set("Authorization", "Bearer "+session.Data.AccessToken)
	router.ServeHTTP(logoutRec, logoutReq)
	if logoutRec.Code != http.StatusNoContent {
		t.Fatalf("logout status = %d, want 204; body: %s", logoutRec.Code, logoutRec.Body.String())
	}

	refreshRec := httptest.NewRecorder()
	refreshReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", jsonBody(`{"refreshToken":"`+session.Data.RefreshToken+`"}`))
	router.ServeHTTP(refreshRec, refreshReq)
	if refreshRec.Code != http.StatusUnauthorized {
		t.Fatalf("refresh status = %d, want 401; body: %s", refreshRec.Code, refreshRec.Body.String())
	}
}

func TestMeReturnsAuthenticatedProfile(t *testing.T) {
	router, _ := newAuthTestRouter(t)

	registerRec := httptest.NewRecorder()
	router.ServeHTTP(registerRec, httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", jsonBody(`{
		"name":"Ada",
		"email":"ada@example.com",
		"password":"correct-password"
	}`)))
	if registerRec.Code != http.StatusCreated {
		t.Fatalf("register status = %d, want 201; body: %s", registerRec.Code, registerRec.Body.String())
	}

	var session dataResponse[authSessionResponse]
	if err := json.Unmarshal(registerRec.Body.Bytes(), &session); err != nil {
		t.Fatalf("decode register response: %v", err)
	}

	meRec := httptest.NewRecorder()
	meReq := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	meReq.Header.Set("Authorization", "Bearer "+session.Data.AccessToken)
	router.ServeHTTP(meRec, meReq)
	if meRec.Code != http.StatusOK {
		t.Fatalf("me status = %d, want 200; body: %s", meRec.Code, meRec.Body.String())
	}

	var profile dataResponse[userResponse]
	if err := json.Unmarshal(meRec.Body.Bytes(), &profile); err != nil {
		t.Fatalf("decode profile: %v", err)
	}
	if profile.Data.Email != "ada@example.com" {
		t.Fatalf("email = %q, want ada@example.com", profile.Data.Email)
	}
}

func jsonBody(body string) *bytes.Reader {
	return bytes.NewReader([]byte(body))
}

func newAuthTestRouter(t *testing.T) (http.Handler, security.TokenManager) {
	t.Helper()

	fixed := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	users := newHTTPMemoryUserRepository(fixed)
	refreshTokens := newHTTPMemoryRefreshTokenRepository(fixed)
	tokens, err := security.NewTokenManager("abcdefghijklmnopqrstuvwxyz1234567890", 15*time.Minute)
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}
	tokens = tokens.WithClock(func() time.Time { return fixed })
	auth, err := service.NewAuthService(service.AuthServiceConfig{
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

	return NewRouter(Dependencies{DB: stubPinger{}, Auth: auth, Tokens: tokens}), tokens
}

type httpMemoryUserRepository struct {
	now     time.Time
	byID    map[uuid.UUID]domain.User
	byEmail map[string]domain.User
}

func newHTTPMemoryUserRepository(now time.Time) *httpMemoryUserRepository {
	return &httpMemoryUserRepository{
		now:     now,
		byID:    map[uuid.UUID]domain.User{},
		byEmail: map[string]domain.User{},
	}
}

func (r *httpMemoryUserRepository) Create(_ context.Context, user domain.User) (domain.User, error) {
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

func (r *httpMemoryUserRepository) FindByID(_ context.Context, id uuid.UUID) (domain.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return domain.User{}, repository.ErrNotFound
	}
	return user, nil
}

func (r *httpMemoryUserRepository) FindByEmail(_ context.Context, email string) (domain.User, error) {
	user, ok := r.byEmail[email]
	if !ok {
		return domain.User{}, repository.ErrNotFound
	}
	return user, nil
}

type httpMemoryRefreshTokenRepository struct {
	now    time.Time
	byID   map[uuid.UUID]domain.RefreshToken
	byHash map[string]domain.RefreshToken
}

func newHTTPMemoryRefreshTokenRepository(now time.Time) *httpMemoryRefreshTokenRepository {
	return &httpMemoryRefreshTokenRepository{
		now:    now,
		byID:   map[uuid.UUID]domain.RefreshToken{},
		byHash: map[string]domain.RefreshToken{},
	}
}

func (r *httpMemoryRefreshTokenRepository) Create(_ context.Context, token domain.RefreshToken) (domain.RefreshToken, error) {
	if _, ok := r.byHash[token.TokenHash]; ok {
		return domain.RefreshToken{}, repository.ErrConflict
	}
	token.ID = uuid.New()
	token.CreatedAt = r.now
	r.byID[token.ID] = token
	r.byHash[token.TokenHash] = token
	return token, nil
}

func (r *httpMemoryRefreshTokenRepository) FindByHash(_ context.Context, tokenHash string) (domain.RefreshToken, error) {
	token, ok := r.byHash[tokenHash]
	if !ok {
		return domain.RefreshToken{}, repository.ErrNotFound
	}
	return token, nil
}

func (r *httpMemoryRefreshTokenRepository) Revoke(_ context.Context, id uuid.UUID) error {
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
