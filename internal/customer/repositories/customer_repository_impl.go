package repository

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
	template "Dzaakk/simple-commerce/package/templates"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
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
	queryCreateCustomer         = `INSERT INTO public.customer (username, email, password, phone_number, balance, created, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	queryFindCustomerById       = `SELECT * FROM public.customer WHERE id = $1`
	queryUpdateBalance          = `UPDATE public.customer SET balance=$1, updated_by=$2, updated=now() WHERE id=$3 RETURNING balance`
	queryUpdateBalanceWithLock  = `UPDATE public.customer SET balance=$1, updated_by='SYSTEM', updated=now() WHERE id=$2 RETURNING balance`
	queryGetBalanceById         = `SELECT id, balance FROM public.customer WHERE id = $1`
	queryGetBalanceByIdWithLock = `SELECT id, balance FROM public.customer WHERE id = $1 FOR UPDATE`
)

func (repo *CustomerRepositoryImpl) Create(data model.TCustomers) (*int, error) {
	statement, err := repo.DB.Prepare(queryCreateCustomer)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var id int
	err = statement.QueryRow(data.Username, data.Email, data.Password, data.PhoneNumber, data.Balance, data.Base.Created, data.Base.CreatedBy).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, err
}

func (repo *CustomerRepositoryImpl) FindById(id int) (*model.TCustomers, error) {
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

func (repo *CustomerRepositoryImpl) FindByEmail(email string) (*model.TCustomers, error) {
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

func (repo *CustomerRepositoryImpl) UpdateBalance(id int, newBalance float64) (*float64, error) {
	statement, err := repo.DB.Prepare(queryUpdateBalance)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var updatedBalance float64
	idString := strconv.Itoa(id)
	err = statement.QueryRow(newBalance, idString, id).Scan(&updatedBalance)
	if err != nil {
		return nil, err
	}

	return &updatedBalance, nil
}

func (repo *CustomerRepositoryImpl) UpdateBalanceWithTx(tx *sql.Tx, id int, newBalance float64) error {
	_, err := tx.Exec(queryUpdateBalanceWithLock, newBalance, id)
	return err
}

func (repo *CustomerRepositoryImpl) GetBalanceWithTx(tx *sql.Tx, id int) (*model.CustomerBalance, error) {
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

func (repo *CustomerRepositoryImpl) GetBalance(id int) (*model.CustomerBalance, error) {
	var customerBalance model.CustomerBalance
	err := repo.DB.QueryRow(queryGetBalanceById, id).Scan(&customerBalance.Id, &customerBalance.Balance)
	if err != nil {
		return nil, err
	}

	return &customerBalance, nil
}

func rowsToCustomer(rows *sql.Rows) (*model.TCustomers, error) {
	base := template.Base{}
	customer := model.TCustomers{}

	err := rows.Scan(&customer.Id, &customer.Username, &customer.Email, &customer.Password, &customer.PhoneNumber, &customer.Balance, &base.Created, &base.CreatedBy, &base.Updated, &base.UpdatedBy)

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
