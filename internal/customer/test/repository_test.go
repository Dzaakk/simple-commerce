package test

import (
	models "Dzaakk/simple-commerce/internal/customer/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCustomer(t *testing.T) {
	repo := NewMockRepository()
	newCustomer := models.TCustomers{
		Username:    "user_test",
		Email:       "test@gmail.com",
		PhoneNumber: "1234567890",
		Password:    "password123",
		Balance:     1000000.00,
	}

	id, err := repo.Create(newCustomer)

	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, 1, *id, "Expected ID to be 1")
}

func TestFindById(t *testing.T) {
	repo := NewMockRepository()

	_, _ = repo.Create(models.TCustomers{
		Username: "user_test",
		Email:    "test@gmail.com",
	})
	customer, err := repo.FindById(1)

	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.Equal(t, "user_test", customer.Username, "Expected username to be user_test")
}

func TestUpdateBalance(t *testing.T) {
	repo := NewMockRepository()

	_, _ = repo.Create(models.TCustomers{
		Username: "user_test",
		Email:    "test@gmail.com",
		Balance:  125000.00,
	})
	updateBalance, err := repo.UpdateBalance(1, 125000.00)

	assert.NoError(t, err)
	assert.NotNil(t, updateBalance)
	assert.Equal(t, float32(125000.00), *updateBalance, "Expected balance to be 125000.00")
}

func TestGetBalance(t *testing.T) {
	repo := NewMockRepository()

	_, _ = repo.Create(models.TCustomers{
		Username: "user_test",
		Email:    "test@gmail.com",
		Balance:  125000,
	})
	customer, err := repo.GetBalance(1)

	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.Equal(t, float32(125000), customer.Balance, "Expected balance to be 125000.00")
}

func TestFindByEmail(t *testing.T) {
	repo := NewMockRepository()

	_, _ = repo.Create(models.TCustomers{
		Username: "user_test",
		Email:    "test@gmail.com",
	})
	customer, err := repo.FindByEmail("test@gmail.com")

	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.Equal(t, "user_test", customer.Username, "Expected username to be user_test")
}
