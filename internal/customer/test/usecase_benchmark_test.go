package test

import (
	models "Dzaakk/simple-commerce/internal/customer/models"
	"Dzaakk/simple-commerce/internal/customer/usecase"
	"testing"
)

func BenchmarkCreateCustomer(b *testing.B) {
	mockRepo := NewMockRepository()
	usecase := usecase.NewCustomerUseCase(mockRepo)

	newCustomer := models.CustomerReq{
		Username:    "user_test",
		Email:       "test@gmail.com",
		PhoneNumber: "123456789",
		Password:    "password123",
	}
	for i := 0; i < b.N; i++ {
		_, err := usecase.Create(newCustomer)
		if err != nil {
			b.Fatalf("Failed to create customer : %v", err)
		}
	}
}
