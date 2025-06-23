package repository

import (
	"Dzaakk/simple-commerce/internal/customer/model"
	response "Dzaakk/simple-commerce/package/response"
	"Dzaakk/simple-commerce/package/template"
	"Dzaakk/simple-commerce/package/util"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
)

var (
	customerTable         = "public.customer"
	customerInsertColumns = []string{
		"username", "email", "password", "gender", "phone_number",
		"balance", "status", "date_of_birth", "profile_picture", "last_login", "created", "created_by",
	}
	customerSelectColumns = append(
		[]string{"id"},
		append(customerInsertColumns, "updated", "updatedby")...,
	)

	QueryFindByEmail              string
	QueryCreate                   string
	QueryFindByID                 string
	QueryUpdateBalance            string
	QueryUpdatePassword           string
	QueryDeactive                 string
	QueryUpdateProfilePic         string
	QueryUpdateBalanceWithReturn  string
	QueryGetBalanceByID           string
	QueryGetBalanceByIDWithReturn string
	once                          sync.Once
)

func InitCustomerQueries() {
	once.Do(func() {
		insertColumns := strings.Join(customerInsertColumns, ",")
		selectColumns := strings.Join(customerSelectColumns, ",")

		insertArgs := util.GeneratePlaceHolders(len(insertColumns))

		QueryCreate = fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s) RRETURNIN id`, customerTable, insertColumns, insertArgs)

		QueryFindByEmail = fmt.Sprintf(`SELECT %s FROM %s WHERE email = $1`, selectColumns, customerTable)
		QueryFindByID = fmt.Sprintf(`SELECT %s FROM %s WHERE id = $1`, selectColumns, customerTable)

		QueryUpdateBalance = `UPDATE public.customer SET balance=$1, updated_by='SYSTEM', updated=now() WHERE id=$2 RETURNING balance`
		QueryUpdatePassword = `UPDATE public.customer SET password=$1, updated_by=$2, updated=now() WHERE id=$2`
		QueryGetBalanceByID = `SELECT balance FROM public.customer WHERE id = $1`
		QueryGetBalanceByIDWithReturn = `SELECT id, balance FROM public.customer WHERE id = $1 FOR UPDATE`
		QueryUpdateProfilePic = "UPDATE public.customer set profile_picture=$1 WHERE id=$2"
		QueryUpdateBalanceWithReturn = `UPDATE public.customer SET balance=$1, updated_by='SYSTEM', updated=now() WHERE id=$2 RETURNING balance`
		QueryDeactive = "UPDATE public.customer set status=$1 WHERE id=$2"
	})
}

type CustomerRepositoryImpl struct {
	DB *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	InitCustomerQueries()
	return &CustomerRepositoryImpl{DB: db}
}

func (repo *CustomerRepositoryImpl) Create(ctx context.Context, data model.TCustomers) (int64, error) {

	result, err := repo.DB.ExecContext(
		ctx, QueryCreate,
		data.Username, data.Email, data.Password,
		data.Gender, data.PhoneNumber, data.Balance,
		data.Status, data.DateOfBirth, data.ProfilePicture,
		data.LastLogin, data.Base.Created, data.Base.CreatedBy,
	)
	if err != nil {
		return 0, response.ExecError("create customer", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, response.Error("failed to retrieve last id", err)
	}

	return id, nil
}

func (repo *CustomerRepositoryImpl) FindByID(ctx context.Context, customerID int64) (*model.TCustomers, error) {
	if customerID <= 0 {
		return nil, response.InvalidParameter()
	}

	row := repo.DB.QueryRowContext(ctx, QueryFindByID, customerID)

	return scanCustomer(row)
}

func (repo *CustomerRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.TCustomers, error) {
	if email == "" {
		return nil, response.InvalidParameter()
	}

	row := repo.DB.QueryRowContext(ctx, QueryFindByEmail, email)

	return scanCustomer(row)
}

func (repo *CustomerRepositoryImpl) UpdateBalance(ctx context.Context, customerID int64, newBalance float64) (int64, error) {
	if customerID <= 0 || newBalance <= -1 {
		return 0, response.InvalidParameter()
	}

	result, err := repo.DB.ExecContext(ctx, QueryUpdateBalance, newBalance, customerID)

	if err != nil {
		return 0, response.ExecError("update balance", err)
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *CustomerRepositoryImpl) UpdatePassword(ctx context.Context, customerID int64, newPassword string) (int64, error) {
	if customerID <= 0 || newPassword == "" {
		return 0, response.InvalidParameter()
	}

	result, err := repo.DB.ExecContext(ctx, QueryUpdatePassword, newPassword, customerID)
	if err != nil {
		return 0, response.ExecError("update password", err)
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *CustomerRepositoryImpl) Deactive(ctx context.Context, customerID int64) (int64, error) {
	if customerID <= 0 {
		return 0, response.InvalidParameter()
	}

	result, err := repo.DB.ExecContext(ctx, QueryDeactive, template.StatusInactive, customerID)
	if err != nil {
		return 0, response.ExecError("deactivate", err)
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *CustomerRepositoryImpl) GetBalance(ctx context.Context, customerID int64) (*model.CustomerBalance, error) {
	if customerID <= 0 {
		return nil, response.InvalidParameter()
	}

	customerBalance := model.CustomerBalance{CustomerID: customerID}

	err := repo.DB.QueryRowContext(ctx, QueryGetBalanceByID, customerID).Scan(&customerBalance.Balance)
	if err != nil {
		return nil, err
	}

	return &customerBalance, nil
}

func (repo *CustomerRepositoryImpl) InquiryBalance(ctx context.Context, customerID int64) (float64, error) {
	if customerID <= 0 {
		return 0, response.InvalidParameter()
	}

	var balance float64

	err := repo.DB.QueryRowContext(ctx, QueryGetBalanceByID, customerID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (repo *CustomerRepositoryImpl) UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, customerID int64, newBalance float64) error {
	if newBalance < 0 || customerID <= 0 {
		return response.InvalidParameter()
	}

	_, err := tx.ExecContext(ctx, QueryUpdateBalance, newBalance, customerID)
	if err != nil {
		return response.ExecError("update with tx", err)
	}

	return nil
}

func (repo *CustomerRepositoryImpl) GetBalanceWithTx(ctx context.Context, tx *sql.Tx, customerID int64) (*model.CustomerBalance, error) {
	row := tx.QueryRowContext(ctx, QueryGetBalanceByIDWithReturn, customerID)

	customerBalance := &model.CustomerBalance{}
	err := row.Scan(&customerBalance.CustomerID, &customerBalance.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, response.Error("customer not found", err)
		}
		return nil, response.Error("error scan customer", err)
	}
	return customerBalance, nil
}

func (repo *CustomerRepositoryImpl) UpdateProfilePicture(ctx context.Context, customerID int64, image string) error {
	if customerID <= 0 {
		return response.InvalidParameter()
	}

	_, err := repo.DB.ExecContext(ctx, QueryUpdateProfilePic, image, customerID)
	if err != nil {
		return response.ExecError("update profile picture", err)
	}

	return nil
}
