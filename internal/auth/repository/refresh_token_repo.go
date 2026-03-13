package repository

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	"Dzaakk/simple-commerce/package/constant"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"errors"
	"time"
)

const (
	refreshTokenSelectColumns        = "id, user_id, user_type, token_hash, expires_at, revoked_at, created_at"
	refreshTokenQueryCreate          = "INSERT INTO public.refresh_tokens (user_id, user_type, token_hash, expires_at, revoked_at, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	refreshTokenQueryFindByUser      = "SELECT " + refreshTokenSelectColumns + " FROM public.refresh_tokens WHERE user_id=$1 ORDER BY created_at DESC LIMIT 1"
	refreshTokenQuerySetExpire       = "UPDATE public.refresh_tokens SET expires_at=$1 WHERE id=$2"
	refreshTokenQueryFindByTokenHash = "SELECT " + refreshTokenSelectColumns + " FROM public.refresh_tokens WHERE token_hash=$1 AND revoked_at IS NULL AND expires_at > NOW() LIMIT 1"
	refreshTokenQueryRevoke          = "UPDATE public.refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1"
	refreshTokenQueryRevokeAllByUser = "UPDATE public.refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND user_type = $2"
)

type RefreshTokenRepository struct {
	DB *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{DB: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, data *model.RefreshToken) (int64, error) {
	var id int64

	err := r.DB.QueryRowContext(
		ctx,
		refreshTokenQueryCreate,
		data.UserID,
		data.UserType,
		data.TokenHash,
		data.ExpiresAt,
		data.RevokedAt,
		data.CreatedAt,
	).Scan(&id)
	if err != nil {
		return 0, response.Error("failed to create refresh token", err)
	}

	return id, nil
}

func (r *RefreshTokenRepository) FindByUserID(ctx context.Context, userID string) (*model.RefreshToken, error) {
	row := r.DB.QueryRowContext(ctx, refreshTokenQueryFindByUser, userID)

	return scanRefreshToken(row)
}

func (r *RefreshTokenRepository) SetExpire(ctx context.Context, id int64, expiresAt time.Time) (int64, error) {
	result, err := r.DB.ExecContext(ctx, refreshTokenQuerySetExpire, expiresAt, id)
	if err != nil {
		return 0, response.ExecError("update refresh token expiry", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, response.Error("failed to get rows affected", err)
	}
	if rowsAffected == 0 {
		return 0, response.Error("no rows updated", sql.ErrNoRows)
	}

	return rowsAffected, nil
}

func (r *RefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	row := r.DB.QueryRowContext(ctx, refreshTokenQueryFindByTokenHash, tokenHash)

	return scanRefreshToken(row)
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	_, err := r.DB.ExecContext(ctx, refreshTokenQueryRevoke, tokenHash)
	if err != nil {
		return response.Error("failed to revoke refresh token", err)
	}

	return nil
}

func (r *RefreshTokenRepository) RevokeAllByUser(ctx context.Context, userID string, userType constant.UserType) error {
	_, err := r.DB.ExecContext(ctx, refreshTokenQueryRevokeAllByUser, userID, userType)
	if err != nil {
		return response.Error("failed to revoke all refresh tokens for user", err)
	}

	return nil
}

func scanRefreshToken(row *sql.Row) (*model.RefreshToken, error) {
	data := &model.RefreshToken{}
	var revokedAt sql.NullTime

	err := row.Scan(
		&data.ID,
		&data.UserID,
		&data.UserType,
		&data.TokenHash,
		&data.ExpiresAt,
		&revokedAt,
		&data.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, response.Error("failed to scan refresh token", err)
	}

	if revokedAt.Valid {
		data.RevokedAt = &revokedAt.Time
	}

	return data, nil
}
