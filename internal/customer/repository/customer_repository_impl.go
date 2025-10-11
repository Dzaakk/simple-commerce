package repository

import (
	"Dzaakk/simple-commerce/internal/customer/model"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
)

const (
	queryCreate      = "INSERT INTO public.customer (username, email, password, gender, phone_number, balance, status, date_of_birth, profile_picture, last_login, created, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id"
	queryFindByEmail = "SELECT id, username, email, password, gender, phone_number, balance, status, date_of_birth, profile_picture, last_login, created, created_by, updated, updated_by FROM public.customer WHERE email=$1"
	queryFindByID    = "SELECT id, username, email, password, gender, phone_number, balance, status, date_of_birth, profile_picture, last_login, created, created_by, updated, updated_by FROM public.customer WHERE id=$1"
	queryUpdate      = "UPDATE public.customer SET username=$1, email=$2, phone_number=$3, date_of_birth=$4, address=$5, updated_by=$6, updated=NOW() WHERE id=$7"
)

type CustomerRepositoryImpl struct {
	DB *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &CustomerRepositoryImpl{DB: db}
}

func (repo *CustomerRepositoryImpl) Create(ctx context.Context, data model.TCustomers) (int64, error) {

	result, err := repo.DB.ExecContext(
		ctx, queryCreate,
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

func (repo *CustomerRepositoryImpl) Update(ctx context.Context, data model.TCustomers) (int64, error) {
	result, err := repo.DB.ExecContext(
		ctx, queryUpdate,
		data.Username, data.Email, data.PhoneNumber,
		data.DateOfBirth, data.Address, data.UpdatedBy,
		data.ID,
	)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (repo *CustomerRepositoryImpl) FindByID(ctx context.Context, customerID int64) (*model.TCustomers, error) {
	if customerID <= 0 {
		return nil, response.InvalidParameter()
	}

	row := repo.DB.QueryRowContext(ctx, queryFindByID, customerID)

	return scanCustomer(row)
}

func (repo *CustomerRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.TCustomers, error) {
	if email == "" {
		return nil, response.InvalidParameter()
	}

	row := repo.DB.QueryRowContext(ctx, queryFindByEmail, email)

	return scanCustomer(row)
}
