package repository

import (
	"Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/constant"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
)

const (
	sellerSelectColumns     = "id, email, password_hash, shop_name, phone, status, created_at, updated_at"
	sellerQueryCreate       = "INSERT INTO public.sellers (id, email, password_hash, shop_name, phone, status, created_at, updated_at) VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7) RETURNING id"
	sellerQueryFindByID     = "SELECT " + sellerSelectColumns + " FROM public.sellers WHERE id=$1"
	sellerQueryFindByEmail  = "SELECT " + sellerSelectColumns + " FROM public.sellers WHERE email=$1"
	sellerQueryFindByName   = "SELECT " + sellerSelectColumns + " FROM public.sellers WHERE shop_name ILIKE $1"
	sellerQueryUpdate       = "UPDATE public.sellers SET email=$1, shop_name=$2, phone=$3, status=$4, updated_at=$5 WHERE id=$6"
	sellerQueryUpdateStatus = "UPDATE public.sellers SET status=$1, updated_at=NOW() WHERE id=$2"
)

type SellerRepository struct {
	DB *sql.DB
}

func NewSellerRepository(db *sql.DB) *SellerRepository {
	return &SellerRepository{DB: db}
}

func (r *SellerRepository) Create(ctx context.Context, data *model.Seller) (string, error) {
	var id string

	err := r.DB.QueryRowContext(
		ctx,
		sellerQueryCreate,
		data.Email,
		data.PasswordHash,
		data.ShopName,
		data.Phone,
		data.Status,
		data.CreatedAt,
		data.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return "", response.Error("failed to create seller", err)
	}

	return id, nil
}

func (r *SellerRepository) Update(ctx context.Context, data *model.Seller) (int64, error) {
	result, err := r.DB.ExecContext(
		ctx,
		sellerQueryUpdate,
		data.Email,
		data.ShopName,
		data.Phone,
		data.Status,
		data.UpdatedAt,
		data.ID,
	)
	if err != nil {
		return 0, response.ExecError("update seller", err)
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

func (r *SellerRepository) FindByID(ctx context.Context, sellerID string) (*model.Seller, error) {
	row := r.DB.QueryRowContext(ctx, sellerQueryFindByID, sellerID)

	return scanSeller(row)
}

func (r *SellerRepository) FindByEmail(ctx context.Context, email string) (*model.Seller, error) {
	row := r.DB.QueryRowContext(ctx, sellerQueryFindByEmail, email)

	return scanSeller(row)
}

func (r *SellerRepository) FindByShopName(ctx context.Context, name string) ([]*model.Seller, error) {
	rows, err := r.DB.QueryContext(ctx, sellerQueryFindByName, "%"+name+"%")
	if err != nil {
		return nil, response.Error("failed to query sellers by shop name", err)
	}
	defer rows.Close()

	var sellers []*model.Seller
	for rows.Next() {
		seller := &model.Seller{}
		if err := rows.Scan(
			&seller.ID,
			&seller.Email,
			&seller.PasswordHash,
			&seller.ShopName,
			&seller.Phone,
			&seller.Status,
			&seller.CreatedAt,
			&seller.UpdatedAt,
		); err != nil {
			return nil, response.Error("failed to scan seller", err)
		}
		sellers = append(sellers, seller)
	}
	if err := rows.Err(); err != nil {
		return nil, response.Error("failed to read sellers", err)
	}

	return sellers, nil
}

func (r *SellerRepository) UpdateStatus(ctx context.Context, sellerID string, status constant.UserStatus) error {
	result, err := r.DB.ExecContext(ctx, sellerQueryUpdateStatus, status, sellerID)
	if err != nil {
		return response.ExecError("update seller status", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return response.Error("failed to get rows affected", err)
	}
	if rowsAffected == 0 {
		return response.Error("no rows updated", sql.ErrNoRows)
	}

	return nil
}

func (r *SellerRepository) UpdateStatusWithTx(ctx context.Context, tx *sql.Tx, sellerID string, status constant.UserStatus) error {
	result, err := tx.ExecContext(ctx, sellerQueryUpdateStatus, status, sellerID)
	if err != nil {
		return response.ExecError("update seller status", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return response.Error("failed to get rows affected", err)
	}
	if rowsAffected == 0 {
		return response.Error("no rows updated", sql.ErrNoRows)
	}

	return nil
}
