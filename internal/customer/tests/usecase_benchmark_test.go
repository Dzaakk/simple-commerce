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

func BenchmarkGetBalanceUsecase(b *testing.B) {
	repo := NewMockRepository()
	usecase := usecase.NewCustomerUseCase(repo)

	_, _ = repo.Create(models.TCustomers{
		Username: "user_test",
		Email:    "test@gmail.com",
		Balance:  125000,
	})

	for i := 0; i < b.N; i++ {
		_, err := usecase.GetBalance(1) // Benchmark the usecase method
		if err != nil {
			b.Errorf("unexpected error: %v", err)
		}
	}
}

func BenchmarkGetBalance(b *testing.B) {
	repo := NewMockRepository()

	_, _ = repo.Create(models.TCustomers{
		Username: "user_test",
		Email:    "test@gmail.com",
		Balance:  125000,
	})

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		_, err := repo.GetBalance(1)
		if err != nil {
			b.Errorf("unexpected error: %v", err)
		}
	}
}

func BenchmarkFindByEmail(b *testing.B) {
	mockRepo := NewMockRepository()

	customer := models.TCustomers{
		Username:    "test_user",
		Email:       "test@example.com",
		PhoneNumber: "123456789",
		Password:    "hashed_password",
		Balance:     100000.0,
	}
	mockRepo.Create(customer)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := mockRepo.FindByEmail("test@example.com")
		if err != nil {
			b.Fatalf("Failed to find customer by email: %v", err)
		}
	}
}

func BenchmarkFindByEmailUseCase(b *testing.B) {
	mockRepo := NewMockRepository()
	usecase := usecase.NewCustomerUseCase(mockRepo)

	customer := models.TCustomers{
		Username:    "test_user",
		Email:       "test@example.com",
		PhoneNumber: "123456789",
		Password:    "hashed_password",
		Balance:     100000.0,
	}
	_, _ = mockRepo.Create(customer)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := usecase.FindByEmail("test@example.com")
		if err != nil {
			b.Fatalf("Failed to find customer by email: %v", err)
		}
	}
}
