package repository

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"errors"
)

const (
	activationCodeSelectColumns         = "id, code, email, type, user_type, expires_at, used_at, created_at"
	activationCodeQueryCreate           = "INSERT INTO public.activation_codes (code, email, type, user_type, expires_at, used_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	activationCodeQueryFindByEmailAndUT = "SELECT " + activationCodeSelectColumns + " FROM public.activation_codes WHERE email=$1 AND user_type=$2 ORDER BY created_at DESC LIMIT 1"
)

type ActivationCodeRepository struct {
	DB *sql.DB
}

func NewActivationCodeRepository(db *sql.DB) *ActivationCodeRepository {
	return &ActivationCodeRepository{DB: db}
}

func (r *ActivationCodeRepository) Create(ctx context.Context, data *model.ActivationCode) (int64, error) {
	var id int64

	err := r.DB.QueryRowContext(
		ctx,
		activationCodeQueryCreate,
		data.Code,
		data.Email,
		data.Type,
		data.UserType,
		data.ExpiresAt,
		data.UsedAt,
		data.CreatedAt,
	).Scan(&id)
	if err != nil {
		return 0, response.Error("failed to create activation code", err)
	}

	return id, nil
}

func (r *ActivationCodeRepository) FindByEmailAndUserType(ctx context.Context, email, userType string) (*model.ActivationCode, error) {
	row := r.DB.QueryRowContext(ctx, activationCodeQueryFindByEmailAndUT, email, userType)

	return scanActivationCode(row)
}

func scanActivationCode(row *sql.Row) (*model.ActivationCode, error) {
	data := &model.ActivationCode{}
	var usedAt sql.NullTime

	err := row.Scan(
		&data.ID,
		&data.Code,
		&data.Email,
		&data.Type,
		&data.UserType,
		&data.ExpiresAt,
		&usedAt,
		&data.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, response.Error("failed to scan activation code", err)
	}

	if usedAt.Valid {
		data.UsedAt = &usedAt.Time
	}

	return data, nil
}
