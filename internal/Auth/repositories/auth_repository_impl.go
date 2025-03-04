package repositories

import (
	"Dzaakk/simple-commerce/internal/auth/models"
	"context"
	"database/sql"
)

type AuthRepositoryImpl struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &AuthRepositoryImpl{DB: db}
}

// CreateCode implements AuthRepository.
func (a *AuthRepositoryImpl) CreateCode(ctx context.Context, data models.TCodeActivation) error {
	panic("unimplemented")
}

// FindCodeByCustomerId implements AuthRepository.
func (a *AuthRepositoryImpl) FindCodeByCustomerId(ctx context.Context, id int64) (*models.TCodeActivation, error) {
	panic("unimplemented")
}
