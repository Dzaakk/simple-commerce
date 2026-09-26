package repository

import (
	"context"
	"database/sql"
	"errors"

	"Dzaakk/simple-commerce/internal/auth/model"
	"Dzaakk/simple-commerce/package/response"
)

const (
	refreshTokenColumns = "id, user_id, token_hash, expires_at, revoked_at, created_at"
	createRefreshToken  = "INSERT INTO refresh_tokens (user_id, token_hash, expires_at, created_at) VALUES ($1, $2, $3, $4) RETURNING id"
	findRefreshToken    = "SELECT " + refreshTokenColumns + " FROM refresh_tokens WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()"
	revokeRefreshToken  = "UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1 AND revoked_at IS NULL"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	return r.db.QueryRowContext(ctx, createRefreshToken,
		token.UserID, token.TokenHash, token.ExpiresAt, token.CreatedAt,
	).Scan(&token.ID)
}

func (r *RefreshTokenRepository) FindByTokenHash(ctx context.Context, hash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	var revokedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, findRefreshToken, hash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &revokedAt, &token.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, response.Error("failed to find refresh token", err)
	}
	if revokedAt.Valid {
		token.RevokedAt = &revokedAt.Time
	}
	return &token, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, hash string) error {
	result, err := r.db.ExecContext(ctx, revokeRefreshToken, hash)
	if err != nil {
		return response.Error("failed to revoke refresh token", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return response.Error("failed to inspect refresh token revocation", err)
	}
	if rows == 0 {
		return response.ErrInvalidRefreshToken
	}
	return nil
}
