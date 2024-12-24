package usecases

import (
	model "Dzaakk/simple-commerce/internal/seller/models"
	repo "Dzaakk/simple-commerce/internal/seller/repositories"
	template "Dzaakk/simple-commerce/package/templates"
	"time"
)

type SellerUseCaseImpl struct {
	repo repo.SellerRepository
}

func NewSellerUseCase(repo repo.SellerRepository) SellerUseCase {
	return &SellerUseCaseImpl{repo}
}

func (s *SellerUseCaseImpl) Create(data model.SellerReq) (*model.TSeller, error) {
	hashedPassword, err := template.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	seller := model.TSeller{
		Username: data.Username,
		Email:    data.Email,
		Password: string(hashedPassword),
		Balance:  0,
		Base: template.Base{
			Created:   time.Now(),
			CreatedBy: "system",
		},
	}

	sellerId, err := 
	
}

// Deactivate implements SellerUseCase.
func (s *SellerUseCaseImpl) Deactivate(sellerId int) (int64, error) {
	panic("unimplemented")
}

// FindById implements SellerUseCase.
func (s *SellerUseCaseImpl) FindById(sellerId int) (*model.TSeller, error) {
	panic("unimplemented")
}

// FindByUsername implements SellerUseCase.
func (s *SellerUseCaseImpl) FindByUsername(username string) (*model.TSeller, error) {
	panic("unimplemented")
}

// Update implements SellerUseCase.
func (s *SellerUseCaseImpl) Update(data model.TSeller) (int64, error) {
	panic("unimplemented")
}
