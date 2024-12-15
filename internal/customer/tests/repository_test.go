package tests

import (
	models "Dzaakk/simple-commerce/internal/customer/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mockRepo     = new(MockCustomerRepository)
	testCustomer = &models.TCustomers{
		Id:          1,
		Username:    "user_test",
		Email:       "test@gmail.com",
		PhoneNumber: "1234567890",
		Password:    "password123",
		Balance:     1000000.00,
	}
	expectedCustomerBalance = &models.CustomerBalance{
		Id:      1,
		Balance: 100000,
	}
	testCustomerID    int
	testCustomerEmail string
	expectedBalance   float64
)

func TestCreateCustomer(t *testing.T) {
	testCustomerID = 1

	mockRepo.On("Create", testCustomer).Return(&testCustomerID, nil)

	createdID, err := mockRepo.Create(testCustomer)

	assert.NoError(t, err)
	assert.NotNil(t, createdID)
	assert.Equal(t, testCustomerID, *createdID)

	mockRepo.AssertExpectations(t)
}

func TestFindById(t *testing.T) {
	testCustomerID = 1

	mockRepo.On("FindById", testCustomerID).Return(testCustomer, nil)

	foundCustomer, err := mockRepo.FindById(testCustomerID)

	assert.NoError(t, err)
	assert.NotNil(t, foundCustomer)

	assert.Equal(t, testCustomer.Username, foundCustomer.Username)
	assert.Equal(t, testCustomer.Email, foundCustomer.Email)
	assert.Equal(t, testCustomer.PhoneNumber, foundCustomer.PhoneNumber)
	assert.Equal(t, testCustomer.Balance, foundCustomer.Balance)

	mockRepo.AssertExpectations(t)
}

func TestUpdateBalance(t *testing.T) {
	testCustomerID = 1
	expectedBalance = 1500000
	mockRepo.On("UpdateBalance", testCustomerID, expectedBalance).Return(&expectedBalance, nil)

	updatedBalance, err := mockRepo.UpdateBalance(testCustomerID, expectedBalance)

	assert.NoError(t, err)
	assert.NotNil(t, updatedBalance)
	assert.Equal(t, &expectedBalance, updatedBalance)

	mockRepo.AssertExpectations(t)
}
func TestGetBalance(t *testing.T) {
	testCustomerID = 1
	mockRepo.On("GetBalance", testCustomerID).Return(expectedCustomerBalance, nil)

	actualCustomerBalance, err := mockRepo.GetBalance(testCustomerID)

	assert.NoError(t, err)
	assert.NotNil(t, actualCustomerBalance)

	assert.Equal(t, expectedCustomerBalance.Id, actualCustomerBalance.Id)
	assert.Equal(t, expectedCustomerBalance.Balance, actualCustomerBalance.Balance)

	mockRepo.AssertExpectations(t)
}

func TestFindByEmail(t *testing.T) {
	testCustomerEmail = "test@gmail.com"
	mockRepo.On("FindByEmail", testCustomerEmail).Return(testCustomer, nil)

	foundCustomer, err := mockRepo.FindByEmail(testCustomerEmail)

	assert.NoError(t, err)
	assert.NotNil(t, foundCustomer)

	assert.Equal(t, testCustomer.Username, foundCustomer.Username)
	assert.Equal(t, testCustomer.Email, foundCustomer.Email)
	assert.Equal(t, testCustomer.PhoneNumber, foundCustomer.PhoneNumber)
	assert.Equal(t, testCustomer.Balance, foundCustomer.Balance)

	mockRepo.AssertExpectations(t)
}

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
