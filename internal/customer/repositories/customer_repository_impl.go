package repositories

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"strconv"
	"time"
)

type CustomerRepositoryImpl struct {
	DB *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &CustomerRepositoryImpl{
		DB: db,
	}
}

const (
	queryFindCustomerByEmail    = `SELECT * FROM public.customer WHERE email = $1`
	queryCreateCustomer         = `INSERT INTO public.customer (username, email, password, phone_number, balance, status, created, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	queryFindCustomerById       = `SELECT * FROM public.customer WHERE id = $1`
	queryUpdateBalance          = `UPDATE public.customer SET balance=$1, updated_by=$2, updated=now() WHERE id=$3 RETURNING balance`
	queryUpdatePassword         = `UPDATE public.customer SET password=$1, updated_by=$2, updated=now() WHERE id=$2`
	queryDeactive               = "UPDATE public.customer set status=$1 WHERE id=$2"
	queryUpdateBalanceWithLock  = `UPDATE public.customer SET balance=$1, updated_by='SYSTEM', updated=now() WHERE id=$2 RETURNING balance`
	queryGetBalanceById         = `SELECT balance FROM public.customer WHERE id = $1`
	queryGetBalanceByIdWithLock = `SELECT id, balance FROM public.customer WHERE id = $1 FOR UPDATE`
	dbQueryTimeout              = 2 * time.Second
)

func (repo *CustomerRepositoryImpl) contextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, dbQueryTimeout)
}

func (repo *CustomerRepositoryImpl) Create(ctx context.Context, data model.TCustomers) (int64, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	result, err := repo.DB.ExecContext(ctx, queryCreateCustomer, data.Username, data.Email, data.Password, data.PhoneNumber, data.Balance, data.Status, data.Base.Created, data.Base.CreatedBy)
	if err != nil {
		return 0, response.ExecError("create customer", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, response.Error("failed to retrieve last id", err)
	}

	return id, err
}

func (repo *CustomerRepositoryImpl) FindById(ctx context.Context, id int64) (*model.TCustomers, error) {
	if id <= 0 {
		return nil, response.InvalidParameter()
	}

	return repo.findCustomer(ctx, queryFindCustomerById, id)
}

func (repo *CustomerRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.TCustomers, error) {
	if email == "" {
		return nil, response.InvalidParameter()
	}

	return repo.findCustomer(ctx, queryFindCustomerByEmail, email)
}

func (repo *CustomerRepositoryImpl) findCustomer(ctx context.Context, query string, args ...interface{}) (*model.TCustomers, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	row := repo.DB.QueryRowContext(ctx, query, args...)
	return scanCustomer(row)
}

func (repo *CustomerRepositoryImpl) UpdateBalance(ctx context.Context, id int64, newBalance float64) (int64, error) {
	if id <= 0 || newBalance <= -1 {
		return 0, response.InvalidParameter()
	}

	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	idString := strconv.FormatInt(id, 8)
	result, err := repo.DB.ExecContext(ctx, queryUpdateBalance, newBalance, idString, id)
	if err != nil {
		return 0, response.ExecError("update balance", err)
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *CustomerRepositoryImpl) UpdatePassword(ctx context.Context, id int64, newPassword string) (int64, error) {
	if id <= 0 || newPassword == "" {
		return 0, response.InvalidParameter()
	}
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	result, err := repo.DB.ExecContext(ctx, queryUpdatePassword, newPassword, id)
	if err != nil {
		return 0, response.ExecError("update password", err)
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *CustomerRepositoryImpl) Deactive(ctx context.Context, id int64) (int64, error) {
	if id <= 0 {
		return 0, response.InvalidParameter()
	}

	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	result, err := repo.DB.ExecContext(ctx, queryDeactive, "I", id)
	if err != nil {
		return 0, response.ExecError("deactivate", err)
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *CustomerRepositoryImpl) GetBalance(ctx context.Context, id int64) (*model.CustomerBalance, error) {
	if id <= 0 {
		return nil, response.InvalidParameter()
	}

	customerBalance := model.CustomerBalance{Id: id}

	err := repo.DB.QueryRow(queryGetBalanceById, id).Scan(&customerBalance.Balance)
	if err != nil {
		return nil, err
	}

	return &customerBalance, nil
}

func (repo *CustomerRepositoryImpl) InquiryBalance(ctx context.Context, id int64) (float64, error) {
	if id <= 0 {
		return 0, response.InvalidParameter()
	}
	row := repo.DB.QueryRowContext(ctx, queryGetBalanceById, id)

	var balance float64
	err := row.Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, response.Error("balance not found", err)
		}
		return 0, response.Error("failed to retrieve balance", err)
	}

	return balance, nil
}

func (repo *CustomerRepositoryImpl) UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, id int64, newBalance float64) error {
	if newBalance < 0 || id <= 0 {
		return response.InvalidParameter()
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, queryUpdateBalanceWithLock, newBalance, id)
	if err != nil {
		return response.ExecError("update with tx", err)
	}

	return nil
}

func (repo *CustomerRepositoryImpl) GetBalanceWithTx(ctx context.Context, tx *sql.Tx, id int64) (*model.CustomerBalance, error) {
	row := tx.QueryRowContext(ctx, queryGetBalanceByIdWithLock, id)

	customerBalance := &model.CustomerBalance{}
	err := row.Scan(&customerBalance.Id, &customerBalance.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, response.Error("customer not found", err)
		}
		return nil, response.Error("error scan customer", err)
	}
	return customerBalance, nil
}
