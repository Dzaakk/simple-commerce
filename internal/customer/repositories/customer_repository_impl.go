package repositories

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
	template "Dzaakk/simple-commerce/package/templates"
	"context"
	"database/sql"
	"errors"
	"fmt"
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
)

func (repo *CustomerRepositoryImpl) Create(ctx context.Context, data model.TCustomers) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	statement, err := repo.DB.PrepareContext(ctx, queryCreateCustomer)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer statement.Close()

	var id int64
	err = statement.QueryRowContext(ctx, data.Username, data.Email, data.Password, data.PhoneNumber, data.Balance, data.Status, data.Base.Created, data.Base.CreatedBy).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return id, err
}

func (repo *CustomerRepositoryImpl) FindById(ctx context.Context, id int64) (*model.TCustomers, error) {
	rows, err := repo.DB.Query(queryFindCustomerById, id)
	if err != nil {
		return nil, err
	}

	customer, err := retrieveCustomer(rows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return customer, nil
}

func (repo *CustomerRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.TCustomers, error) {
	rows, err := repo.DB.Query(queryFindCustomerByEmail, email)
	if err != nil {
		return nil, err
	}

	customer, err := retrieveCustomer(rows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return customer, nil
}

func (repo *CustomerRepositoryImpl) UpdateBalance(ctx context.Context, id int64, newBalance float64) (float64, error) {
	statement, err := repo.DB.Prepare(queryUpdateBalance)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	var updatedBalance float64
	idString := strconv.FormatInt(id, 8)
	err = statement.QueryRow(newBalance, idString, id).Scan(&updatedBalance)
	if err != nil {
		return 0, err
	}

	return updatedBalance, nil
}

func (repo *CustomerRepositoryImpl) UpdatePassword(ctx context.Context, id int64, newPassword string) (int64, error) {
	statement, err := repo.DB.Prepare(queryUpdatePassword)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(newPassword, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *CustomerRepositoryImpl) Deactive(ctx context.Context, id int64) (int64, error) {
	statement, err := repo.DB.Prepare(queryDeactive)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec("I", id)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *CustomerRepositoryImpl) UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, id int64, newBalance float64) error {
	_, err := tx.Exec(queryUpdateBalanceWithLock, newBalance, id)
	return err
}

func (repo *CustomerRepositoryImpl) GetBalanceWithTx(ctx context.Context, tx *sql.Tx, id int64) (*model.CustomerBalance, error) {
	var customerBalance model.CustomerBalance
	err := tx.QueryRow(queryGetBalanceByIdWithLock, id).
		Scan(&customerBalance.Id, &customerBalance.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer with id %d not found", id)
		}
		return nil, fmt.Errorf("unable to retrieve customer balance : %w", err)
	}
	return &customerBalance, nil
}

func (repo *CustomerRepositoryImpl) GetBalance(ctx context.Context, id int64) (*model.CustomerBalance, error) {
	customerBalance := model.CustomerBalance{Id: id}

	err := repo.DB.QueryRow(queryGetBalanceById, id).Scan(&customerBalance.Balance)
	if err != nil {
		return nil, err
	}

	return &customerBalance, nil
}

func (repo *CustomerRepositoryImpl) InquiryBalance(ctx context.Context, id int64) (float64, error) {
	var balance float64
	err := repo.DB.QueryRow(queryGetBalanceById, id).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("invalid customer id")
		}
		return 0, fmt.Errorf("failed to retrieve balance")
	}

	return balance, nil
}

func rowsToCustomer(rows *sql.Rows) (*model.TCustomers, error) {
	base := template.Base{}
	customer := model.TCustomers{}

	err := rows.Scan(&customer.Id, &customer.Username, &customer.Email, &customer.Password, &customer.PhoneNumber, &customer.Balance, &customer.Status, &base.Created, &base.CreatedBy, &base.Updated, &base.UpdatedBy)

	if err != nil {
		return nil, err
	}
	if !base.UpdatedBy.Valid {
		base.UpdatedBy.String = ""
	}
	customer.Base = base

	return &customer, nil
}
func retrieveCustomer(rows *sql.Rows) (*model.TCustomers, error) {
	if rows.Next() {
		return rowsToCustomer(rows)
	}
	return nil, errors.New("customer not found")
}
