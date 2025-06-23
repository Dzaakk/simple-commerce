package repository

import (
	"Dzaakk/simple-commerce/internal/seller/model"
	"Dzaakk/simple-commerce/package/response"
	"Dzaakk/simple-commerce/package/template"
	"Dzaakk/simple-commerce/package/util"
	"context"
	"fmt"
	"strings"
	"sync"

	"database/sql"
	"time"
)

var (
	sellerTable         = "public.seller"
	sellerInsertColumns = []string{
		"username", "email", "password", "phone_number", "store_name", "address", "balance", "status", "bank_account_name", "bank_account_number", "bank_name", "created", "created_by",
	}
	sellerSelectColumns = append(
		[]string{"id"},
		append(sellerInsertColumns, "updated", "updated_by")...,
	)

	QueryCreate         string
	QueryUpdate         string
	QueryUpdatePassword string
	QueryDeactive       string
	QueryFindBySellerID string
	QueryFindAll        string
	QueryFindByUsername string
	QueryFindByEmail    string
	QueryUpdateBalance  string
	once                sync.Once
)

func InitSellerQueries() {
	once.Do(func() {
		insertColumns := strings.Join(sellerInsertColumns, ",")
		selectColumns := strings.Join(sellerSelectColumns, ",")

		insertArgs := util.GeneratePlaceHolders(len(insertColumns))

		QueryCreate = fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s) RRETURNIN id`, sellerTable, insertColumns, insertArgs)

		QueryFindAll = fmt.Sprintf(`SELECT %s FROM %s ORDER BY username`, selectColumns, sellerTable)
		QueryFindBySellerID = fmt.Sprintf(`SELECT %s FROM %s WHERE id = $1`, selectColumns, sellerTable)
		QueryFindByUsername = fmt.Sprintf(`SELECT %s FROM %s WHERE username = $1`, selectColumns, sellerTable)
		QueryFindByEmail = fmt.Sprintf(`SELECT %s FROM %s WHERE email = $1`, selectColumns, sellerTable)

		QueryUpdate = "UPDATE public.seller SET username=$1, email=$2, updated=NOW(), updated_by=$3 WHERE id=$4"
		QueryUpdatePassword = "UPDATE public.seller set password=$1 WHERE id=$2"
		QueryDeactive = "UPDATE public.seller set status=$1 WHERE id=$2"
		QueryUpdateBalance = "UPDATE public.seller SET balance=$1, updated=NOW(), updated_by=$2 WHERE id=$2"
	})
}

type SellerRepositoryImpl struct {
	DB *sql.DB
}

func NewSellerRepository(db *sql.DB) SellerRepository {
	InitSellerQueries()
	return &SellerRepositoryImpl{DB: db}
}

func (repo *SellerRepositoryImpl) Create(ctx context.Context, data model.TSeller) (int64, error) {

	result, err := repo.DB.ExecContext(
		ctx, QueryCreate,
		data.Username, data.Email, data.Password, data.PhoneNumber,
		data.StoreName, data.Address, data.Balance, data.Status,
		data.BankAccountName, data.BankAccountNumber, data.BankName,
		data.Created, data.CreatedBy,
	)
	if err != nil {
		return 0, response.ExecError("create seller", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, response.Error("failed to retrieve last id", err)
	}

	return id, nil
}

func (repo *SellerRepositoryImpl) Update(ctx context.Context, data model.TSeller) (int64, error) {
	result, err := repo.DB.ExecContext(ctx, QueryUpdate, data.Username, data.Email, time.Now(), data.Username)
	if err != nil {
		return 0, err
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *SellerRepositoryImpl) FindAll(ctx context.Context) ([]*model.TSeller, error) {
	rows, err := repo.DB.QueryContext(ctx, QueryFindAll)
	if err != nil {
		return nil, response.Error("error scan seller", err)
	}
	defer rows.Close()

	return scanListSeller(rows)
}

func (repo *SellerRepositoryImpl) FindBySellerID(ctx context.Context, sellerID int64) (*model.TSeller, error) {
	row := repo.DB.QueryRowContext(ctx, QueryFindBySellerID, sellerID)

	return scanSeller(row)
}
func (repo *SellerRepositoryImpl) FindByUsername(ctx context.Context, username string) (*model.TSeller, error) {
	row := repo.DB.QueryRowContext(ctx, QueryFindByUsername, username)

	return scanSeller(row)
}
func (repo *SellerRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.TSeller, error) {
	row := repo.DB.QueryRowContext(ctx, QueryFindByEmail, email)

	return scanSeller(row)
}

func (repo *SellerRepositoryImpl) InsertBalance(ctx context.Context, sellerID int64, balance int64) error {
	_, err := repo.DB.ExecContext(ctx, QueryUpdateBalance, balance, sellerID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SellerRepositoryImpl) UpdatePassword(ctx context.Context, sellerID int64, newPassword string) (int64, error) {

	result, err := repo.DB.ExecContext(ctx, QueryUpdatePassword, newPassword, sellerID)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *SellerRepositoryImpl) Deactive(ctx context.Context, sellerID int64) (int64, error) {

	result, err := repo.DB.ExecContext(ctx, QueryDeactive, template.StatusInactive, sellerID)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}
