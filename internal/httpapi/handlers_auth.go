package httpapi

import (
	"net/http"
	"time"

	"github.com/WassimBannout/gatherops/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) AuthHandler {
	return AuthHandler{auth: auth}
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type logoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type dataResponse[T any] struct {
	Data T `json:"data"`
}

type authSessionResponse struct {
	User                  userResponse `json:"user"`
	AccessToken           string       `json:"accessToken"`
	AccessTokenExpiresAt  time.Time    `json:"accessTokenExpiresAt"`
	RefreshToken          string       `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time    `json:"refreshTokenExpiresAt"`
}

type userResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (h AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error(), nil)
		return
	}

	session, err := h.auth.Register(r.Context(), service.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		WriteServiceError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusCreated, dataResponse[authSessionResponse]{Data: authSessionToResponse(session)})
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error(), nil)
		return
	}

	session, err := h.auth.Login(r.Context(), service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		WriteServiceError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, dataResponse[authSessionResponse]{Data: authSessionToResponse(session)})
}

func (h AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error(), nil)
		return
	}

	session, err := h.auth.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		WriteServiceError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, dataResponse[authSessionResponse]{Data: authSessionToResponse(session)})
}

func (h AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := AuthenticatedUserID(r.Context())
	if !ok {
		WriteError(w, r, http.StatusUnauthorized, "unauthorized", "Authentication required", nil)
		return
	}

	var req logoutRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error(), nil)
		return
	}

	if err := h.auth.Logout(r.Context(), userID, req.RefreshToken); err != nil {
		WriteServiceError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := AuthenticatedUserID(r.Context())
	if !ok {
		WriteError(w, r, http.StatusUnauthorized, "unauthorized", "Authentication required", nil)
		return
	}

	profile, err := h.auth.Me(r.Context(), userID)
	if err != nil {
		WriteServiceError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, dataResponse[userResponse]{Data: userToResponse(profile)})
}

func authSessionToResponse(session service.AuthSession) authSessionResponse {
	return authSessionResponse{
		User:                  userToResponse(session.User),
		AccessToken:           session.AccessToken,
		AccessTokenExpiresAt:  session.AccessTokenExpiresAt,
		RefreshToken:          session.RefreshToken,
		RefreshTokenExpiresAt: session.RefreshTokenExpiresAt,
	}
}

func userToResponse(user service.UserProfile) userResponse {
	return userResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
