package repositories

import (
	"Dzaakk/simple-commerce/internal/auth/models"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"time"
)

type AuthRepositoryImpl struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &AuthRepositoryImpl{DB: db}
}

const (
	queryCreate    = `INSERT INTO public.code_activation (user_id, code_activation, is_used, created_at, used_at) VALUES ($1, $2, $3, $4, $5)`
	dbQueryTimeout = 3 * time.Second
)

func (repo *AuthRepositoryImpl) contextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, dbQueryTimeout)
}

func (repo *AuthRepositoryImpl) CreateCode(c context.Context, data models.TCodeActivation) error {
	ctx, cancel := repo.contextWithTimeout(c)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, queryCreate, data.UserID, data.CodeActivation, data.IsUsed, data.CreatedAt, data.UsedAt)
	if err != nil {
		return response.ExecError("create activation code", err)
	}

	return nil
}

// FindCodeByCustomerId implements AuthRepository.
func (a *AuthRepositoryImpl) FindCodeByCustomerId(ctx context.Context, id int64) (*models.TCodeActivation, error) {
	panic("unimplemented")
}
