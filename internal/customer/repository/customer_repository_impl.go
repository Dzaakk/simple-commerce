package repositories

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"log"
	"time"
)

const (
	StatusActive   = "A"
	StatusInactive = "I"
	updatedBy      = "SYSTEM"
)

const (
	queryFindByEmail              = `SELECT * FROM public.customer WHERE email = $1`
	queryCreate                   = `INSERT INTO public.customer (username, email, password, phone_number, balance, status, created, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	queryFindByID                 = `SELECT * FROM public.customer WHERE id = $1`
	queryUpdateBalance            = `UPDATE public.customer SET balance=$1, updated_by=$2, updated=now() WHERE id=$3 RETURNING balance`
	queryUpdatePassword           = `UPDATE public.customer SET password=$1, updated_by=$2, updated=now() WHERE id=$2`
	queryDeactive                 = "UPDATE public.customer set status=$1 WHERE id=$2"
	queryUpdateBalanceWithReturn  = `UPDATE public.customer SET balance=$1, updated_by='SYSTEM', updated=now() WHERE id=$2 RETURNING balance`
	queryGetBalanceByID           = `SELECT balance FROM public.customer WHERE id = $1`
	queryGetBalanceByIDWithReturn = `SELECT id, balance FROM public.customer WHERE id = $1 FOR UPDATE`
	dbQueryTimeout                = 2 * time.Second
)

type CustomerRepositoryImpl struct {
	DB    *sql.DB
	stmts map[string]*sql.Stmt
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	repo := &CustomerRepositoryImpl{
		DB:    db,
		stmts: make(map[string]*sql.Stmt),
	}

	prepareQueries := map[string]string{
		"findByCustomerID":    queryFindByID,
		"findByCustomerEmail": queryFindByEmail,
		"updateBalance":       queryUpdateBalance,
		// "updateBalanceWithReturn": queryUpdateBalanceByCustomerIDWithReturn,
		"getBalance": queryGetBalanceByID,
		// "getBalanceWithReturn":    queryGetBalanceByCustomerIDWithReturn,
		"updatePassword": queryUpdatePassword,
	}

	for key, query := range prepareQueries {
		stmt, err := db.Prepare(query)
		if err != nil {
			repo.Close()
			log.Printf("failed to prepare %s statement: %v", key, err)
			return nil
		}
		repo.stmts[key] = stmt
	}

	return repo
}

// close all prepared statements
func (repo *CustomerRepositoryImpl) Close() error {
	for name, stmt := range repo.stmts {
		if err := stmt.Close(); err != nil {
			log.Printf("Error closing statement %s: %v", name, err)
		}
	}
	return nil
}

func (repo *CustomerRepositoryImpl) contextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining < dbQueryTimeout {
			return context.WithCancel(ctx)
		}
	}
	return context.WithTimeout(ctx, dbQueryTimeout)
}

func (repo *CustomerRepositoryImpl) Create(ctx context.Context, data model.TCustomers) (int64, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	result, err := repo.DB.ExecContext(
		ctx,
		queryCreate,
		data.Username,
		data.Email,
		data.Password,
		data.PhoneNumber,
		data.Balance,
		data.Status,
		data.Base.Created,
		data.Base.CreatedBy)
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

	return repo.findCustomer(ctx, queryFindByID, customerID)
}

func (repo *CustomerRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.TCustomers, error) {
	if email == "" {
		return nil, response.InvalidParameter()
	}

	return repo.findCustomer(ctx, queryFindByEmail, email)
}

func (repo *CustomerRepositoryImpl) findCustomer(ctx context.Context, query string, args ...interface{}) (*model.TCustomers, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	var row *sql.Row
	if query == queryFindByID && repo.stmts["findByCustomerID"] != nil {
		row = repo.stmts["findByCustomerID"].QueryRowContext(ctx, args...)
	} else if query == queryFindByEmail && repo.stmts["findByCustomerEmail"] != nil {
		row = repo.stmts["findByCustomerEmail"].QueryRowContext(ctx, args...)
	} else {
		row = repo.DB.QueryRowContext(ctx, query, args...)
	}
	return scanCustomer(row)
}

func (repo *CustomerRepositoryImpl) UpdateBalance(ctx context.Context, customerID int64, newBalance float64) (int64, error) {
	if customerID <= 0 || newBalance <= -1 {
		return 0, response.InvalidParameter()
	}

	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	var result sql.Result
	var err error

	if repo.stmts["updateBalance"] != nil {
		result, err = repo.stmts["updateBalance"].ExecContext(ctx, newBalance, updatedBy, customerID)
	} else {
		result, err = repo.DB.ExecContext(ctx, queryUpdateBalance, newBalance, updatedBy, customerID)
	}

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
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	result, err := repo.DB.ExecContext(ctx, queryUpdatePassword, newPassword, customerID)
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

	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	result, err := repo.DB.ExecContext(ctx, queryDeactive, StatusInactive, customerID)
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

	var err error
	if repo.stmts["getBalance"] != nil {
		err = repo.stmts["getBalance"].QueryRowContext(ctx, customerID).Scan(&customerBalance.Balance)
	} else {
		err = repo.DB.QueryRowContext(ctx, queryGetBalanceByID, customerID).Scan(&customerBalance.Balance)
	}
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
	var err error

	if repo.stmts["getBalance"] != nil {
		err = repo.stmts["getBalance"].QueryRowContext(ctx, customerID).Scan(&balance)
	} else {
		err = repo.DB.QueryRowContext(ctx, queryGetBalanceByID, customerID).Scan(&balance)
	}
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (repo *CustomerRepositoryImpl) UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, customerID int64, newBalance float64) error {
	if newBalance < 0 || customerID <= 0 {
		return response.InvalidParameter()
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, queryUpdateBalance, newBalance, customerID)
	if err != nil {
		return response.ExecError("update with tx", err)
	}

	return nil
}

func (repo *CustomerRepositoryImpl) GetBalanceWithTx(ctx context.Context, tx *sql.Tx, customerID int64) (*model.CustomerBalance, error) {
	row := tx.QueryRowContext(ctx, queryGetBalanceByIDWithReturn, customerID)

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
