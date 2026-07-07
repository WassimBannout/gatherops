package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/WassimBannout/gatherops/internal/domain"
	"github.com/WassimBannout/gatherops/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepository(pool *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{pool: pool}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token domain.RefreshToken) (domain.RefreshToken, error) {
	const query = `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, token_hash, expires_at, revoked_at, created_at
	`

	created, err := scanRefreshToken(r.pool.QueryRow(ctx, query, token.UserID, token.TokenHash, token.ExpiresAt))
	if err != nil {
		if isUniqueViolation(err) {
			return domain.RefreshToken{}, repository.ErrConflict
		}
		return domain.RefreshToken{}, fmt.Errorf("create refresh token: %w", err)
	}
	return created, nil
}

func (r *RefreshTokenRepository) FindByHash(ctx context.Context, tokenHash string) (domain.RefreshToken, error) {
	const query = `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	token, err := scanRefreshToken(r.pool.QueryRow(ctx, query, tokenHash))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.RefreshToken{}, repository.ErrNotFound
		}
		return domain.RefreshToken{}, fmt.Errorf("find refresh token by hash: %w", err)
	}
	return token, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	const query = `
		UPDATE refresh_tokens
		SET revoked_at = COALESCE(revoked_at, now())
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	if result.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

type refreshTokenRow interface {
	Scan(dest ...any) error
}

func scanRefreshToken(row refreshTokenRow) (domain.RefreshToken, error) {
	var token domain.RefreshToken
	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
	)
	return token, err
}
