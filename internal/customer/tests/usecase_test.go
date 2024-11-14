package test

import (
	models "Dzaakk/simple-commerce/internal/customer/models"
	u "Dzaakk/simple-commerce/internal/customer/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomerUsecase_Create(t *testing.T) {
	mockRepo := NewMockRepository()
	usecase := u.NewCustomerUseCase(mockRepo)

	newCustomer := models.CustomerReq{
		Username:    "user_test",
		Email:       "test@gmail.com",
		PhoneNumber: "123456789",
		Password:    "password123",
	}

	id, err := usecase.Create(newCustomer)

	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, 1, *id)
}

func TestCustomerUsecase_GetBalance(t *testing.T) {
	mockRepo := NewMockRepository()
	usecase := u.NewCustomerUseCase(mockRepo)

	_, _ = mockRepo.Create(models.TCustomers{
		Username: "user_test",
		Email:    "test@gmail.com",
		Balance:  125000,
	})

	customer, err := usecase.GetBalance(1)

	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.Equal(t, float32(125000), customer.Balance)
}

func TestCustomerUsecase_UpdateBalance(t *testing.T) {
	mockRepo := NewMockRepository()
	usecase := u.NewCustomerUseCase(mockRepo)

	_, _ = mockRepo.Create(models.TCustomers{
		Username: "user_test",
		Email:    "test@gmail.com",
		Balance:  25000,
	})

	newBalance, err := usecase.UpdateBalance(1, 100000, "A")

	assert.NoError(t, err)
	assert.NotNil(t, newBalance)
	assert.Equal(t, float32(125000), *newBalance)
}
