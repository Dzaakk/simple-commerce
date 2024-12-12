package tests

import (
	models "Dzaakk/simple-commerce/internal/customer/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mockRepo     = new(MockCustomerRepository)
	CustomerMock = &models.TCustomers{
		Id:          1,
		Username:    "user_test",
		Email:       "test@gmail.com",
		PhoneNumber: "1234567890",
		Password:    "password123",
		Balance:     1000000.00,
	}
	expectedID int
)

func TestCreateCustomer(t *testing.T) {
	expectedID = 1

	mockRepo.On("Create", CustomerMock).Return(&expectedID, nil)

	createdID, err := mockRepo.Create(CustomerMock)

	assert.NoError(t, err)
	assert.NotNil(t, createdID)
	assert.Equal(t, expectedID, *createdID)

	mockRepo.AssertExpectations(t)
}

func TestFindById(t *testing.T) {
	expectedID = 1

	mockRepo.On("FindById", expectedID).Return(CustomerMock, nil)

	foundCustomer, err := mockRepo.FindById(expectedID)

	assert.NoError(t, err)
	assert.NotNil(t, foundCustomer)

	assert.Equal(t, CustomerMock.Username, foundCustomer.Username)
	assert.Equal(t, CustomerMock.Email, foundCustomer.Email)
	assert.Equal(t, CustomerMock.PhoneNumber, foundCustomer.PhoneNumber)
	assert.Equal(t, CustomerMock.Balance, foundCustomer.Balance)

	// Verify that the expected method was called
	mockRepo.AssertExpectations(t)
}

// func TestUpdateBalance(t *testing.T) {
// 	repo := NewMockRepository()

// 	_, _ = repo.Create(models.TCustomers{
// 		Username: "user_test",
// 		Email:    "test@gmail.com",
// 		Balance:  125000.00,
// 	})
// 	updateBalance, err := repo.UpdateBalance(1, 125000.00)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, updateBalance)
// 	assert.Equal(t, float32(125000.00), *updateBalance, "Expected balance to be 125000.00")
// }

// func TestGetBalance(t *testing.T) {
// 	repo := NewMockRepository()

// 	_, _ = repo.Create(models.TCustomers{
// 		Username: "user_test",
// 		Email:    "test@gmail.com",
// 		Balance:  125000,
// 	})
// 	customer, err := repo.GetBalance(1)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, customer)
// 	assert.Equal(t, float32(125000), customer.Balance, "Expected balance to be 125000.00")
// }

// func TestFindByEmail(t *testing.T) {
// 	repo := NewMockRepository()

// 	_, _ = repo.Create(models.TCustomers{
// 		Username: "user_test",
// 		Email:    "test@gmail.com",
// 	})
// 	customer, err := repo.FindByEmail("test@gmail.com")

// 	assert.NoError(t, err)
// 	assert.NotNil(t, customer)
// 	assert.Equal(t, "user_test", customer.Username, "Expected username to be user_test")
// }
